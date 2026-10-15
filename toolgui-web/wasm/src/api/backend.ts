import { AppConf, DownloadResult, UpdateEvent, UploadResult } from "@toolgui-web/lib"

type DirectoryPickerWindow = Window & {
  showDirectoryPicker?: (options: { mode: 'read' | 'readwrite' }) => Promise<FileSystemDirectoryHandle>
}

// LoadProgress is how far the download of app.wasm has got. total is 0 when
// the size is unknown; done is set once the last byte is in.
export interface LoadProgress {
  loaded: number
  total: number
  done?: boolean
}

// Backend drives the wasm program in its worker. Calls go out as messages and
// come back by id; packs arrive on their own, the way the websocket transport
// delivers them.
export class Backend {
  // Leave the adapter absent in browsers that do not support directory grants.
  // This arrow is passed directly through FolderPickerContext.
  readonly pickFolder = typeof (window as DirectoryPickerWindow).showDirectoryPicker === 'function'
    ? (componentID: string, _defaultDir: string, write: boolean) => this.pickDirectory(componentID, write)
    : undefined
  private worker: Worker
  private pending = new Map<number, (msg: any) => void>()
  private nextID = 1
  private generation = 0
  private ready: Promise<void>

  // onPack is called for every pack, in the order the page produced them.
  // embed is the display mode, which the program is told at boot: a worker's
  // own location is this script, so it cannot read the page's query string.
  // query is the page's query string, handed over for the same reason.
  // onProgress follows the download of the binary, until the program starts.
  constructor(onPack: (pack: any) => void, embed: boolean = false,
    onProgress?: (progress: LoadProgress) => void, query: string = '') {
    this.worker = new Worker(new URL('../worker.ts', import.meta.url))

    let started: () => void
    let failed: (e: any) => void
    this.ready = new Promise((resolve, reject) => { started = resolve; failed = reject })

    this.worker.onmessage = (event: MessageEvent) => {
      const msg = event.data

      switch (msg.kind) {
        case 'ready':
          started()
          break

        case 'failed':
          failed(new Error(msg.error))
          break

        case 'pack':
          onPack(msg.pack)
          break

        case 'event':
          // Prefixed as tgwasm.EventPrefix, so it cannot pass for a browser event.
          window.dispatchEvent(new CustomEvent('toolgui:' + msg.name, { detail: msg.detail }))
          break

        case 'progress':
          onProgress?.({ loaded: msg.loaded, total: msg.total, done: msg.done })
          break

        case 'return': {
          const resolve = this.pending.get(msg.id)
          this.pending.delete(msg.id)
          resolve?.(msg)
          break
        }

        default:
          console.error('unknown message', msg.kind)
      }
    }

    // Both files sit next to index.html, wherever the app is hosted.
    this.worker.postMessage({
      kind: 'init',
      wasmExecURL: assetURL('wasm_exec.js'),
      wasmURL: assetURL('app.wasm'),
      embed,
      query,
      base: new URL('.', document.baseURI).href,
    })
  }

  // send posts one message the worker answers by id, whatever it is.
  private async send(msg: any): Promise<any> {
    await this.ready

    const id = this.nextID++
    return new Promise((resolve, reject) => {
      this.pending.set(id, (answer) => {
        if (answer.error) {
          reject(new Error(answer.error))
          return
        }

        resolve(answer.value)
      })

      try {
        this.worker.postMessage({ ...msg, id })
      } catch (e) {
        this.pending.delete(id)
        reject(e)
      }
    })
  }

  private call(fn: string, ...args: any[]): Promise<any> {
    return this.send({ kind: 'call', fn, args })
  }

  private async pickDirectory(componentID: string, write: boolean): Promise<void> {
    const generation = this.generation
    // Do not await the worker before opening the dialog: transient activation
    // belongs to this click, on the main thread. Reserve the originating session
    // in parallel so navigation during the dialog cannot hand a grant to another.
    const reservation = this.call('beginFolderPick', componentID, write)
    let picker: Promise<FileSystemDirectoryHandle>
    try {
      picker = (window as DirectoryPickerWindow).showDirectoryPicker({ mode: write ? 'readwrite' : 'read' })
    } catch (e) {
      picker = Promise.reject(e)
    }
    const [slot, selected] = await Promise.allSettled([reservation, picker])
    if (slot.status === 'rejected') throw slot.reason
    if (generation !== this.generation) {
      await this.call('cancelFolderPick', slot.value)
      return
    }
    if (selected.status === 'rejected') {
      await this.call('cancelFolderPick', slot.value)
      if (['AbortError', 'NotAllowedError'].includes(selected.reason?.name)) return
      throw selected.reason
    }
    try {
      // Directory handles are structured-cloneable; their bytes and paths do
      // not cross the bridge. The worker calls the Go bridge with this handle.
      await this.call('folderPicked', slot.value, selected.value)
    } finally {
      await this.call('cancelFolderPick', slot.value)
    }
  }

  // appConf is the browser counterpart of GET /api/app.
  async appConf(): Promise<AppConf> {
    return JSON.parse(await this.call('appConf'))
  }

  // start opens a session on a page and draws it once. query is the page
  // query, `group=a`.
  start(pageName: string, query: string): Promise<void> {
    this.generation++
    return this.call('start', pageName, query)
  }

  update(event: UpdateEvent): Promise<void> {
    return this.call('update', JSON.stringify(event))
  }

  // uploadFile is the browser counterpart of POST /api/files.
  //
  // The File itself crosses, not its content. A structured clone of one hands
  // the worker the same blob this thread holds -- the browser has it on disk
  // already -- so nothing is copied here and nothing this size is ever a
  // string. The worker streams it into the file system the wasm program reads
  // uploads from; see worker.ts.
  async uploadFile(file: File, componentID: string): Promise<UploadResult> {
    try {
      await this.send({ kind: 'upload', componentID, file })
      return { ok: true }
    } catch (e) {
      return { ok: false, error: String(e) }
    }
  }

  // downloadFile is the browser counterpart of GET /api/files.
  //
  // Nothing crosses but the token and, coming back, the file the worker read
  // out of the origin private file system: a structured clone of it carries
  // the blob's backing rather than its bytes, so a download costs this thread
  // nothing whatever it weighs.
  async downloadFile(token: string): Promise<DownloadResult> {
    try {
      const file = await this.send({ kind: 'download', token })
      return { ok: true, blob: file }
    } catch (e) {
      return { ok: false, error: String(e) }
    }
  }
}

function assetURL(name: string): string {
  return new URL(name, document.baseURI).href
}
