import { useState } from "react"
import { sessionValues } from "./session"

export function TCheckbox({ node, update }) {
  const [value, setValue] = useState(sessionValues[node.props.id] || false)
  return (
    <div class="field">
      <div class="control">
        <label class="checkbox">
          <input type="checkbox"
            id={node.props.id}
            checked={value}
            onChange={(event) => {
              sessionValues[event.target.id] = event.target.checked
              setValue(event.target.checked)
              update({
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
