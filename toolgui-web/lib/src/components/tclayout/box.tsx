import React from "react"

import { Props } from "../component_interface"
import { TComponent } from "../factory"

export function TBox({ node, update, nodes }: Props) {
  return (
    <div id={node.props.id} className="box">
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
