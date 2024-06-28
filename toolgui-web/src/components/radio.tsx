import { useState } from "react"
import { sessionValues } from "./session"
import { Props } from "./component_interface"

export function TRadio({ node, update }: Props) {
  const items: string[] = node.props.items
  const [value, setValue] = useState<string>(sessionValues[node.props.id] || '')
  return (
    <div className="control">
      {
        items.map(x =>
          <label className="radio">
            <input type="radio" name={node.props.id} value={x}
              onChange={(event) => {
                sessionValues[node.props.id] = event.target.value
                setValue(event.target.value)
                update({
                  id: node.props.id,
                  value: event.target.value,
                })
              }}
              checked={value === x} />
            {x}
          </label>
        )
      }
    </div>
  )
}
