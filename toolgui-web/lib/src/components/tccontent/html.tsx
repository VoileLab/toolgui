import React from "react"

import { Props } from '../component_interface'

export function THtml({ node }: Props) {
  const sandbox = node.props.script ? "allow-scripts allow-same-origin" : "allow-same-origin"

  return (
    <iframe
      id={node.props.id}
      name={node.props.id}
      sandbox={sandbox}
      srcDoc={node.props.html}
    />
  )
}
