import React, { useState } from "react"
import { stateValues } from "../state"
import { Props } from "../component_interface"

export function TDatepicker({ node, update }: Props) {
  const [value, setValue] = useState<string>(stateValues[node.props.id] || '')
  return (
    <div className="field">
      <label className="label">{node.props.label}</label>
      <div className="control">
        <input type={node.props.type}
          className="input"
          id={node.props.id}
          value={value}
          onChange={(event) => {
            stateValues[event.target.id] = event.target.value
            setValue(event.target.value)
          }}
          onBlur={(event) => {
            update({
              id: event.target.id,
              value: stateValues[event.target.id],
            })
          }}>
        </input>
      </div>
    </div>
  )
}
