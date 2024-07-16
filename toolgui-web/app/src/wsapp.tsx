import React, { Component } from "react"
import { initHealthSock, uploadFile, wsUpdate } from "./api/updater"

import { App } from "@toolgui-web/lib"
import { AppConf } from "@toolgui-web/lib"
import { StatefulWebSocket } from "./api/StatefulWebSocket"

interface WSState {
  appConf: AppConf | null
  pageName: string | null
  conn: StatefulWebSocket | null
}

export class WSApp extends Component<{}, WSState> {
  appEle: React.RefObject<App>

  constructor(props: any) {
    super(props)
    this.state = {
      appConf: null,
      pageName: null,
      conn: null,
    }
    this.getAppConf()

    this.appEle = React.createRef()
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

      const conn = new StatefulWebSocket(pageName, pack => {
        if (pack.success !== undefined) {
          this.appEle.current.finishUpdate(pack)
          return
        }

        this.appEle.current.receiveNotifyPack(pack)

      }, () => {
        this.appEle.current.clearState()
      })

      this.setState({ appConf, pageName, conn })
    }).catch(e => {
      console.log(e)
    })
  }

  update(event: any) {
    this.appEle.current.startUpdate()
    this.state.conn.send(event)
  }

  render(): React.ReactNode {
    if (!this.state.appConf) {
      return <></>
    }

    return (
      <App appConf={this.state.appConf}
        ref={this.appEle}
        update={(pack) => { this.update(pack) }}
        upload={uploadFile} />
    )
  }
}