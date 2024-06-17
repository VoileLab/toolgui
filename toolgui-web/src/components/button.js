export function TButton({ node, update }) {
  return (
    <button id={node.props.id} class="button" onClick={
      (event) => {
        update({
          id: event.target.id,
          value: true,
          // revoke state change after running finish
          is_temp: true,
        })
      }
    }>
      {node.props.label}
    </button>
  )
}
