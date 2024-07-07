import { Component } from "react"
import { initHealthSock, wsUpdate } from "./api/updater"

import { App } from "@toolgui-web/lib"
import { AppConf } from "@toolgui-web/lib"

interface WSState {
  appConf: AppConf | null
  pageName: string | null
}

export class WSApp extends Component<{}, WSState> {
  constructor(props: any) {
    super(props)
    this.state = {
      appConf: null,
      pageName: null,
    }
    this.getAppConf()
  }

  getAppConf() {
    fetch('/api/app').then(resp => resp.json()).then((appConf: AppConf) => {
      var pageName = ''
      if (appConf.hash_page_name_mode) {
        if (window.location.hash) {
          // should be #/{name}
          pageName = window.location.hash.substring(2)
        } else if (appConf.page_names.length > 0) {
          pageName = appConf.page_names[0]
        }
      } else {
        pageName = window.location.pathname.substring(1)
      }

      initHealthSock(pageName)

      this.setState({ appConf, pageName })
    }).catch(e => {
      console.log(e)
    })
  }

  update(
    event: any,
    clearContainer: () => void,
    clearState: () => void,
    recvNotifyPack: (pack: any) => void,
    finishUpdate: (pack: any) => void) {

    wsUpdate(this.state.pageName, event, clearContainer, clearState, recvNotifyPack, finishUpdate)
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