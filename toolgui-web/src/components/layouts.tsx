import { Props } from "./component_interface"
import { TComponent } from "./factory"

export function TContainer({ node, update, nodes }: Props) {
  return (
    <div id={node.props.id}>
      {
        node.children.map(child =>
          <TComponent node={child}
            update={update}
            nodes={nodes} />
        )
      }
    </div>
  )
}

export function TBox({ node, update, nodes }: Props) {
  return (
    <div id={node.props.id} className="box">
      {
        node.children.map(child =>
          <TComponent node={child}
            update={update}
            nodes={nodes} />
        )
      }
    </div>
  )
}

export function TColumn({ node, update, nodes }: Props) {
  return (
    <div id={node.props.id} className="columns">
      {
        node.children.map(child =>
          <div className="column">
            <TComponent node={child}
              update={update}
              nodes={nodes} />
          </div>
        )
      }
    </div>
  )
}
