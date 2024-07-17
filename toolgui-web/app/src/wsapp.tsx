import React, { Component } from "react"

import { App } from "@toolgui-web/lib"
import { AppConf } from "@toolgui-web/lib"
import { StatefulWebSocket } from "./api/StatefulWebSocket"
import { getAppConf } from "./api/AppConfAPI"

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
    this.appEle = React.createRef()

    this.setup()
  }

  async setup() {
    const appConf = await getAppConf()

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

    const conn = new StatefulWebSocket(pageName, pack => {
      if (pack.success !== undefined) {
        if (!pack.success) {
          console.error(pack)
        }

        this.appEle.current.finishUpdate(pack)
        return
      }

      if (pack.ready !== undefined) {
        if (!pack.ready) {
          console.error(pack)
        }

        this.appEle.current.startUpdate()
        return
      }

      this.appEle.current.receiveNotifyPack(pack)
    })

    conn.onConnect = () => {
      this.state.conn.send({})
    }

    conn.onStateIDChange = () => {
      this.appEle.current.clearState()
    }

    conn.init()

    this.setState({ appConf, pageName, conn })
  }

  render(): React.ReactNode {
    if (!this.state.appConf) {
      return <></>
    }

    return (
      <App appConf={this.state.appConf}
        ref={this.appEle}
        update={(pack) => { this.state.conn.send(pack) }}
        upload={(f) => { return this.state.conn.uploadFile(f) }} />
    )
  }
}