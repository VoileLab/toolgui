import React from "react"

import { stateValues } from "../state"
import { Props } from "../component_interface"

export function TCheckbox({ node, update }: Props) {
  return (
    <div className="field">
      <div className="control">
        <label className="checkbox">
          <input type="checkbox"
            id={node.props.id}
            checked={stateValues[node.props.id] || node.props.default}
            disabled={node.props.disabled}
            onChange={(event) => {
              stateValues[event.target.id] = event.target.checked
              update({
                type: "input",
                id: event.target.id,
                value: event.target.checked,
              })
            }} />
          &nbsp;{node.props.label}
        </label>
      </div>
    </div>
  )
}
