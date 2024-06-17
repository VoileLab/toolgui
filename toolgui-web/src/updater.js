var sessionID = ''
var sock = null

export function updater(event, clearContainer, createComponent) {
    const pageName = window.location.pathname.substring(1)

    if (sessionID != '') {
        event['session_id'] = sessionID
    }

    if (sock) {
        sock.close()
    }

    sock = new WebSocket(`/api/update/${pageName}`)
    var jsonEvent = JSON.stringify(event)

    sock.onopen = function () {
        clearContainer()
        sock.send(jsonEvent)
    }

    sock.onmessage = function (e) {
        const data = JSON.parse(e.data)
        if (data.session_id) {
            sessionID = data.session_id
            return
        }

        createComponent(data)
    }
}
