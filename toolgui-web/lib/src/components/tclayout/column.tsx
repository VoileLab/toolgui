import React from "react"

import { Props } from "../component_interface"
import { TComponent } from "../factory"

export function TColumn({ node, update, upload }: Props) {
  return (
    <div id={node.props.id} className="columns">
      {
        node.children.map(child =>
          <div className="column">
            <TComponent node={child}
              update={update}
              upload={upload} />
          </div>
        )
      }
    </div>
  )
}
