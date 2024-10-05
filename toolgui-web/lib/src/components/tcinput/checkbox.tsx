import React from "react"

import { stateValues } from "../state"
import { Props } from "../component_interface"

export function TCheckbox({ node, update }: Props) {
  return (
    <div className="field">
      <div className="control">
        <label className="label">
          <input type="checkbox"
            id={node.props.id}
            checked={stateValues[node.props.id]}
            onChange={(event) => {
              stateValues[event.target.id] = event.target.checked
              update({
                type: "input",
                id: event.target.id,
                value: event.target.checked,
              })
            }} />
          {node.props.label}
        </label>
      </div>
    </div>
  )
}
