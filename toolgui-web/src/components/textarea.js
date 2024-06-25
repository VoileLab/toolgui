import { useState } from "react"
import { sessionValues } from "./session"

export function TTextarea({ node, update }) {
  const [value, setValue] = useState(sessionValues[node.props.id])
  return (
    <div class="field">
      <label class="label">{node.props.label}</label>
      <div class="control">
        <textarea class="textarea"
          id={node.props.id}
          value={value}
          onChange={(event) => {
            sessionValues[event.target.id] = event.target.value
            setValue(event.target.value)
          }}
          onBlur={(event) => {
            update({
              id: event.target.id,
              value: sessionValues[event.target.id],
            })
          }}>
        </textarea>
      </div>
    </div>
  )
}
