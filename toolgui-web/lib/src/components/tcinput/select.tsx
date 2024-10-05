import React from "react";

import { stateValues } from "../state"
import { Props } from "../component_interface"

export function TSelect({ node, update }: Props) {
  var value = ''
  if (stateValues[node.props.id] > 0) {
    value = node.props.items[stateValues[node.props.id] - 1]
  }

  return (
    <div className="field">
      <label className="label">{node.props.label}</label>
      <div className="select">
        <select
          id={node.props.id}
          value={value}
          onChange={(event) => {
            const index = event.target.selectedIndex
            stateValues[event.target.id] = index
            update({
              type: "input",
              id: event.target.id,
              value: index,
            })
          }}>
          <option value="">Please select an option</option>
          {
            node.props.items.map((item: string) =>
              <option value={item}>{item}</option>)
          }
        </select>
      </div>
    </div>
  )
}
