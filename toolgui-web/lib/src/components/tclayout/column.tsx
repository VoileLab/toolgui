import React from "react"

import { Props } from "../component_interface"
import { TComponent } from "../factory"

export function TColumn({ node, nodes, update, upload }: Props) {
  return (
    <div id={node.props.id} className="columns">
      {
        node.children.map(child =>
          <div className="column">
            <TComponent node={child}
              nodes={nodes}
              update={update}
              upload={upload} />
          </div>
        )
      }
    </div>
  )
}
