import { Props } from "./component_interface";

export function TProgressar({ node }: Props) {
  return (
    <div>
      <p>{node.props.label}</p>
      <progress id={node.props.id}
        className="progress is-primary"
        value={node.props.value} max="100">
      </progress>
    </div>
  )
}
