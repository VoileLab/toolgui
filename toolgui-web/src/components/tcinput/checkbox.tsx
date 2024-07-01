import { sessionValues } from "../session"
import { Props } from "../component_interface"

export function TCheckbox({ node, update }: Props) {
  return (
    <div className="field">
      <div className="control">
        <label className="label">
          <input type="checkbox"
            id={node.props.id}
            checked={sessionValues[node.props.id]}
            onChange={(event) => {
              sessionValues[event.target.id] = event.target.checked
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
