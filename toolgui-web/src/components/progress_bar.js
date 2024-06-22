export function TProgressar({ node }) {
  return (
    <div>
      <p>{node.props.label}</p>
      <progress id={node.props.id}
        class="progress is-primary"
        value={node.props.value} max="100">
      </progress>
    </div>
  )
}
