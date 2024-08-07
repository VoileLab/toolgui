import React from "react";

import { stateValues } from "../state"
import { Props } from "../component_interface"

export function TSelect({ node, update }: Props) {
  return (
    <div className="field">
      <label className="label">{node.props.label}</label>
      <div className="select">
        <select
          id={node.props.id}
          value={stateValues[node.props.id]}
          onChange={(event) => {
            stateValues[event.target.id] = event.target.value
            update({
              id: event.target.id,
              value: event.target.value,
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
