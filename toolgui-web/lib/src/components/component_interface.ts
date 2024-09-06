import { Node } from "../app/Nodes"
import { UpdateEvent } from "../app/UpdateEvent"

export interface Props {
  node: Node

  update: (event: UpdateEvent) => void
  upload: (file: File) => Promise<Response>
}
