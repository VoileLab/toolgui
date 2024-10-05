import React from "react";

import { stateValues } from "../state"
import { Props } from "../component_interface"

export function TRadio({ node, update }: Props) {
  const items: string[] = node.props.items
  return (
    <div className="control">
      {
        items.map((x, idx) =>
          <label className="radio" key={idx}>
            <input type="radio" name={node.props.id} value={idx}
              onChange={(event) => {
                stateValues[node.props.id] = Number(event.target.value)
                update({
                  type: "input",
                  id: node.props.id,
                  value: Number(event.target.value),
                })
              }}
              checked={stateValues[node.props.id] === idx} />
            {x}
          </label>
        )
      }
    </div>
  )
}
