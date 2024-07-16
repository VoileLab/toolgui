const healthCheck = "/api/health"

enum WebSocketState {
  Initial,
  Ping,
  TryConnect,
  Connected,

  // TODO: Dead
}

enum WebSocketAction {
  OK,
  Error,
  Closed,

  // TODO: support timeout check
}

function sleep(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

function getSocketURI() {
  var scheme = 'ws'
  if (window.location.origin.startsWith('https')) {
    scheme = 'wss'
  }

  return `${scheme}://${window.location.host}`
}

function getUpdateURI(pageName: string) {
  return `${getSocketURI()}/api/update/${pageName}`
}

export class StatefulWebSocket {
  conn: WebSocket | null
  state: WebSocketState
  pageName: string
  stateID: string
  recv: (pack: any) => void
  clearState: () => void

  constructor(pageName: string, recv: (pack: any) => void, clearState: () => void) {
    this.state = WebSocketState.Initial
    this.pageName = pageName
    this.conn = null
    this.stateID = ''
    this.recv = recv
    this.clearState = clearState
    this.walk(WebSocketAction.OK)
  }

  walkTo(state: WebSocketState) {
    switch (state) {
      case WebSocketState.Ping:
        this.state = state
        this.ping()
        break
      case WebSocketState.TryConnect:
        this.state = state
        this.tryConnect()
        break
      case WebSocketState.Connected:
        this.state = state
        break
      default:
        console.error('undefined state', state)
        throw new Error('undefine state')
    }
  }

  walk(action: WebSocketAction) {
    console.log('current state', this.state)
    console.log('walk', action)
    switch (this.state) {
      case WebSocketState.Initial:
        if (action === WebSocketAction.OK) {
          this.walkTo(WebSocketState.Ping)
          return
        }
        return
      case WebSocketState.Ping:
        if (action == WebSocketAction.OK) {
          this.walkTo(WebSocketState.TryConnect)
          return
        }
      case WebSocketState.TryConnect:
        switch (action) {
          case WebSocketAction.OK:
            this.walkTo(WebSocketState.Connected)
            return
          case WebSocketAction.Error:
          case WebSocketAction.Closed:
            this.walkTo(WebSocketState.Ping)
            return
          default:
            break
        }
        break
      case WebSocketState.Connected:
        switch (action) {
          case WebSocketAction.Error:
          case WebSocketAction.Closed:
            this.walkTo(WebSocketState.Ping)
            return
          default:
            break
        }
        break
    }

    console.error('undefine action on state', this.state, action)
    throw new Error('undefine action on state')
  }

  async ping() {
    while (1) {
      const resp = await fetch(healthCheck)
      if (resp.status === 200) {
        break
      }
      console.error('health check', resp)
      await sleep(1000)
    }
    console.log('ping ok')
    this.walk(WebSocketAction.OK)
  }

  tryConnect() {
    this.conn = new WebSocket(getUpdateURI(this.pageName))
    var that = this

    this.conn.onopen = function () {
      that.conn.send(JSON.stringify({ state_id: that.stateID }))
      console.log('socket open ok')
      that.walk(WebSocketAction.OK)
    }

    this.conn.onmessage = function (e) {
      const data = JSON.parse(e.data)
      if (data.state_id) {
        that.stateID = data.state_id
        that.clearState()
        return
      }

      that.recv(data)
    }

    this.conn.onclose = function () {
      that.conn = null
      that.walk(WebSocketAction.Closed)
    }
  }

  send(pack: any) {
    if (this.state !== WebSocketState.Connected || this.conn === null) {
      throw new Error('websocket is not prepared')
    }

    this.conn.send(JSON.stringify(pack))
  }
}
