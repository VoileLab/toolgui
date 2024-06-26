var sessionID = ''

var updateSock = null
var healthSock = null

function getSocketURI() {
    var scheme = 'ws'
    if (window.location.origin.startsWith('https')) {
        scheme = 'wss'
    }

    return `${scheme}://${window.location.host}`
}

function getUpdateURI() {
    const pageName = window.location.pathname.substring(1)
    return `${getSocketURI()}/api/update/${pageName}`
}

function getHealthURI() {
    const pageName = window.location.pathname.substring(1)
    return `${getSocketURI()}/api/health/${pageName}`
}

export function updater(event,
    clearContainer, clearSession, recvNotifyPack, finishUpdate) {

    if (sessionID !== '') {
        event['session_id'] = sessionID
    }

    if (updateSock) {
        updateSock.close()
    }

    updateSock = new WebSocket(getUpdateURI())
    var jsonEvent = JSON.stringify(event)

    updateSock.onopen = function () {
        clearContainer()
        updateSock.send(jsonEvent)
    }

    updateSock.onmessage = function (e) {
        const data = JSON.parse(e.data)
        if (data.session_id) {
            sessionID = data.session_id
            clearSession()
            return
        }

        if (data.success !== undefined) {
            finishUpdate(data)
            return
        }

        recvNotifyPack(data)
    }
}

export function initHealthSock() {
    if (healthSock) {
        return
    }

    healthSock = new WebSocket(getHealthURI())
    healthSock.onopen = function () {
        console.log('Start health beating')
    }

    // health beat / 1 mins
    setInterval(function () {
        healthSock.send(JSON.stringify({ session_id: sessionID }))
    }, 60 * 1000);
}