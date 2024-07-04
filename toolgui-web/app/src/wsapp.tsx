import { Component } from "react"
import { initHealthSock, wsUpdate } from "./api/updater"

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
    this.getAppConf()
  }

  componentDidMount() {
    initHealthSock()
  }

  getAppConf() {
    fetch('/api/app').then(resp => resp.json()).then((appConf: AppConf) => {
      this.setState({ appConf })
    }).catch(e => {
      console.log(e)
    })
  }

  render(): React.ReactNode {
    return (
      <>
        {this.state.appConf ?
          <App appConf={this.state.appConf} updater={wsUpdate} /> : ''}
      </>
    )
  }
}