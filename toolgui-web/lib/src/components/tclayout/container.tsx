import React from "react"

import { Props } from "../component_interface"
import { TComponent } from "../factory"

export function TContainer({ node, nodes, update, upload }: Props) {
  return (
    <div id={node.props.id}>
      {
        node.children.map(child =>
          <TComponent node={child}
            nodes={nodes}
            update={update}
            upload={upload} />
        )
      }
    </div>
  )
}