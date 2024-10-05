export interface ClickEvent {
  type: "click"
  id: string
}

export interface InputEvent {
  type: "input"
  id: string
  value: any
}

export type UpdateEvent = ClickEvent | InputEvent
