import React from "react"

import { Props } from "../component_interface"
import { TComponent } from "../factory"

export function TExpand({ node, update, upload, theme }: Props) {
  return (
    <details id={node.props.id} className="box" open={node.props.expanded}>
      <summary>
        {node.props.title}
      </summary>
      {
        node.children.map(child =>
          <TComponent node={child}
            update={update}
            upload={upload}
            theme={theme} />
        )
      }
    </details>
  )
}