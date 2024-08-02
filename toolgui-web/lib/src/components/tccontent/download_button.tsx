import React from "react"

import { Props } from '../component_interface'

export function TDownloadButton({ node }: Props) {
  return (
    <a id={node.props.id}
      href={'data:application/octet-stream;base64,' + node.props.base64_body}
      download={node.props.filename}>
      {node.props.text}
    </a>
  )
}
