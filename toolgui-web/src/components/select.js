import { useState } from "react"
import { sessionValues } from "./session"

export function TSelect({ node, update }) {
  const [value, setValue] = useState(sessionValues[node.props.id] || '')
  return (
    <div class="select">
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
          node.props.items.map((item) =>
            <option value={item}>{item}</option>)
        }
      </select>
    </div>
  )
}
