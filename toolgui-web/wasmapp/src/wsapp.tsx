import { Component } from "react"

import { App } from "@toolgui-web/lib"
import { AppConf } from "@toolgui-web/lib"

interface WSState {
  appConf: AppConf | null
}

export class WSApp extends Component<{}, WSState> {
  constructor(props: any) {
    super(props)
    this.state = {
      appConf: null
    }
    this.initGo()
  }

  initGo() {
    const go = new Go()
    WebAssembly.instantiateStreaming(fetch('main.wasm'), go.importObject).then(
      ({ instance }) => {
        go.run(instance)
        const appConf: AppConf = JSON.parse(getAppConf())
        this.setState({ appConf })
      }
    )
  }

  update(
    event: any,
    clearContainer: () => void,
    clearSession: () => void,
    recvNotifyPack: (pack: any) => void,
    finishUpdate: (pack: any) => void) {

  }

  render(): React.ReactNode {
    return (
      <>
        {this.state.appConf ?
          <App appConf={this.state.appConf} updater={(a, b, c, d, e) => { this.update(a, b, c, d, e) }} /> : ''}
      </>
    )
  }
}