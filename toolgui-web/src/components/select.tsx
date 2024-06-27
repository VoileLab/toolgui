import { useState } from "react"
import { sessionValues } from "./session"
import { Props } from "./component_interface"

export function TSelect({ node, update }: Props) {
  const [value, setValue] = useState<string>(sessionValues[node.props.id] || '')
  return (
    <div className="select">
      <select
        id={node.props.id}
        value={value}
        onChange={(event) => {
          sessionValues[event.target.id] = event.target.value
          setValue(event.target.value)
          update({
            id: event.target.id,
            value: event.target.value,
          })
        }}>
        {
          node.props.items.map((item: string) =>
            <option value={item}>{item}</option>)
        }
      </select>
    </div>
  )
}
