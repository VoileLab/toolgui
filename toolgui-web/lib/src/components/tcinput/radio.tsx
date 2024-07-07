import React from "react";

import { stateValues } from "../state"
import { Props } from "../component_interface"

export function TRadio({ node, update }: Props) {
  const items: string[] = node.props.items
  return (
    <div className="control">
      {
        items.map(x =>
          <label className="radio">
            <input type="radio" name={node.props.id} value={x}
              onChange={(event) => {
                stateValues[node.props.id] = event.target.value
                update({
                  id: node.props.id,
                  value: event.target.value,
                })
              }}
              checked={stateValues[node.props.id] === x} />
            {x}
          </label>
        )
      }
    </div>
  )
}
