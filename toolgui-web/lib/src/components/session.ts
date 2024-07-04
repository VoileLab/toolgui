export var sessionValues: { [key: string]: any } = {}

export function clearSession() {
    sessionValues = {}
}
