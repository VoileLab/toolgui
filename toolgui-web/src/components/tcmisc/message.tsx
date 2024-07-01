import { Props } from "../component_interface";

export function TMessage({ node }: Props) {
  return (
    <article className={`message is-${node.props.type}`}>
      <div className="message-header">
        <p>{node.props.title}</p>
      </div>
      <div className="message-body">
        {node.props.text}
      </div>
    </article>
  )
}

