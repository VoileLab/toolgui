export interface ClickEvent {
  type: "click"
  id: string
}

export interface InputEvent {
  type: "input"
  id: string
  value: any
}

export interface SelectEvent {
  type: "select"
  id: string
  value: number
}

export type UpdateEvent = ClickEvent | InputEvent | SelectEvent
