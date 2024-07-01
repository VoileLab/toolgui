import { sessionValues } from "../session"
import { Props } from "../component_interface"

export function TDatepicker({ node, update }: Props) {
  return (
    <div className="field">
      <label className="label">{node.props.label}</label>
      <div className="control">
        <input type="date"
          className="input"
          id={node.props.id}
          value={sessionValues[node.props.id]}
          onChange={(event) => {
            sessionValues[event.target.id] = event.target.value
            update({
              id: event.target.id,
              value: sessionValues[event.target.id],
            })
          }}>
        </input>
      </div>
    </div>
  )
}
