import React from "react"

import { Props } from "../component_interface"
import { TComponent } from "../factory"

export function TContainer({ node, update, nodes }: Props) {
  return (
    <div id={node.props.id}>
      {
        node.children.map(child =>
          <TComponent node={child}
            update={update}
            nodes={nodes} />
        )
      }
    </div>
  )
}