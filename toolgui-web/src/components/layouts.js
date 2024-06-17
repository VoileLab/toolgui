import { TComponent } from "./factory"

export function TContainer({ node, update, nodes }) {
  return (
    <div id={node.props.id}>
      {
        node.children.map(name =>
          <TComponent node={nodes[name]}
            update={update}
            nodes={nodes} />
        )
      }
    </div>
  )
}

export function TBox({ node, update, nodes }) {
  return (
    <div id={node.props.id} class="box">
      {
        node.children.map(name =>
          <TComponent node={nodes[name]}
            update={update}
            nodes={nodes} />
        )
      }
    </div>
  )
}

export function TColumn({ node, update, nodes }) {
  return (
    <div id={node.props.id} class="columns">
      {
        node.children.map(name =>
          <div class="column">
            <TComponent node={nodes[name]}
              update={update}
              nodes={nodes} />
          </div>
        )
      }
    </div>
  )
}
