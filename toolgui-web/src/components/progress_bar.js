export function TProgressar({ node }) {
  return (
    <progress id={node.props.id} class="progress" value={node.props.value} max="100">
      { node.props.label }
    </progress>
  )
}
