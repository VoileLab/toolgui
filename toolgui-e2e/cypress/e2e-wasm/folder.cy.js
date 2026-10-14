// A real OPFS directory handle is structured-cloned to the actual Go worker.
// Only the native picker is stubbed; file IO and the page reruns are real.
describe('Wasm folder picker', () => {
  let directory, picker, root, name
  beforeEach(() => {
    cy.visit('/#/folder', {
      onBeforeLoad(win) {
        name = `folder-e2e-${Date.now()}-${Math.random()}`
        const fixture = win.navigator.storage.getDirectory().then(async origin => {
          root = origin
          directory = await origin.getDirectoryHandle(name, { create: true })
          const file = await directory.getFileHandle('input.txt', { create: true })
          const stream = await file.createWritable()
          await stream.write('hello')
          await stream.close()
          return directory
        })
        picker = cy.stub().callsFake(() => fixture)
        win.showDirectoryPicker = picker
      },
    })
    cy.contains('button', 'Choose local folder').should('be.enabled')
  })
  afterEach(() => {
    cy.then(async () => { if (root) await root.removeEntry(name, { recursive: true }) })
  })

  it('grants read access, reads and lists files, then keeps state on cancellation and denial', () => {
    cy.contains('button', 'Choose local folder').click()
    cy.contains('Selected folder: read-only').should('exist')
    cy.then(() => expect(picker.firstCall.args[0]).to.deep.equal({ mode: 'read' }))
    cy.contains('Entry: input.txt').should('exist')
    cy.contains('button', 'Read local file').click()
    cy.contains('Local content: hello').should('exist')
    cy.contains('File size: 5').should('exist')
    cy.contains('button', 'Write local file').should('not.exist')
    for (const error of ['AbortError', 'NotAllowedError']) {
      cy.window().then(win => picker.callsFake(() => Promise.reject(new win.DOMException('cancelled', error))))
      cy.contains('button', 'Choose local folder').click().should('be.enabled')
      // The result of the read button disappears on rerun, so this also checks
      // that a rejected selection did not send any state event.
      cy.contains('Local content: hello').should('exist')
      cy.get('[role=alert]').should('not.exist')
    }
  })

  it('requests readwrite and uses the same grant for directory creation and file replacement', () => {
    cy.contains('label', 'Request write access').click()
    cy.contains('button', 'Choose local folder').click()
    cy.contains('Selected folder: readwrite').should('exist')
    cy.then(() => expect(picker.firstCall.args[0]).to.deep.equal({ mode: 'readwrite' }))
    cy.contains('button', 'Create directories').click()
    cy.contains('Directories created').should('exist')
    cy.then(() => directory.getDirectoryHandle('nested').then(dir => dir.getDirectoryHandle('deep')))
    cy.contains('button', 'Write local file').click()
    cy.contains('Local file written').should('exist')
    cy.then(() => directory.getFileHandle('input.txt').then(file => file.getFile()).then(file => file.text()))
      .should('equal', 'written by Go\n')
    cy.contains('button', 'Read local file').click()
    cy.contains('Local content: written by Go').should('exist')
  })

  it('discards a selection finished after navigation to a replacement session', () => {
    let finish
    cy.window().then(win => picker.callsFake(() => new win.Promise(resolve => { finish = resolve })))
    cy.contains('button', 'Choose local folder').click()
    cy.window().then(win => { win.location.hash = '#/index' })
    cy.contains('button', 'Choose local folder').should('not.exist')
    cy.window().then(win => { win.location.hash = '#/folder' })
    cy.contains('button', 'Choose local folder').should('be.enabled')
    cy.then(() => finish(directory))
    cy.contains('Selected folder:').should('not.exist')
  })
})

describe('Wasm without directory picker support', () => {
  it('leaves folder selection disabled', () => {
    cy.visit('/#/folder', { onBeforeLoad(win) { win.showDirectoryPicker = undefined } })
    cy.contains('button', 'Choose local folder').should('be.disabled')
    cy.contains('Folder selection is not available in this executor yet.').should('exist')
  })
})
