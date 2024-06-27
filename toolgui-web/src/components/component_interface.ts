import { Node } from "../Nodes"

export interface Event {
  session_id?: string
  id?: string
  value?: any

  // revoke state change after running finish
  is_temp?: boolean
}

export interface Props {
  node: Node
  update: (event: Event) => void
  nodes: { [id: string]: Node }
}
