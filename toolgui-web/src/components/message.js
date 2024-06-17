export function TMessage({ node }) {
  return (
    <article class={`message is-${node.props.type}`}>
      <div class="message-header">
        <p>{node.props.title}</p>
      </div>
      <div class="message-body">
        {node.props.text}
      </div>
    </article>
  )
}

