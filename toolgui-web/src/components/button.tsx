import { Props } from "./component_interface"

export function TButton({ node, update }: Props) {
  return (
    <button id={node.props.id} className="button" onClick={
      (event: React.MouseEvent<HTMLButtonElement>) => {
        const target = event.target as HTMLButtonElement
        update({
          id: target.id,
          value: true,
          is_temp: true,
        })
      }
    }>
      {node.props.label}
    </button>
  )
}
