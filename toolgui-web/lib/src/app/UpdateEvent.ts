export interface UpdateEvent {
  state_id?: string
  id?: string
  value?: any

  // revoke state change after running finish
  is_temp?: boolean
}
