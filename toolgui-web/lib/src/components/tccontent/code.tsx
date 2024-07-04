import React from 'react'

import hljs from 'highlight.js'
import 'highlight.js/styles/default.css'
import { Props } from '../component_interface'

export function TCode({ node }: Props) {
  const highlightHTML = hljs.highlight(
    node.props.code,
    { language: node.props.lang }
  ).value
  return (
    <div className="content">
      <pre>
        <div dangerouslySetInnerHTML={{ __html: highlightHTML }} />
      </pre>
    </div>
  )
}
