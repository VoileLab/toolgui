import { useState } from "react"
import { sessionValues } from "./session"
import { Props } from "./component_interface"

export function TCheckbox({ node, update }: Props) {
  const [value, setValue] = useState<boolean>(sessionValues[node.props.id] || false)
  return (
    <div className="field">
      <div className="control">
        <label className="label">
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
