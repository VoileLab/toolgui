import React, { useState } from "react"
import { stateValues } from "../state"
import { Props } from "../component_interface"

export function TTextbox({ node, update }: Props) {
  const [value, setValue] = useState<string>(stateValues[node.props.id] || '')

  var inputClassNames = 'input'
  if (node.props.color !== '') {
    inputClassNames += ' is-' + node.props.color
  }

  return (
    <div className="field">
      <label className="label">{node.props.label}</label>
      <div className="control">
        <input
          type={node.props.password ? 'password' : 'text'}
          maxLength={node.props.max_length ? node.props.max_length : ''}
          placeholder={node.props.placeholder}
          disabled={node.props.disabled}
          className={inputClassNames}
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
