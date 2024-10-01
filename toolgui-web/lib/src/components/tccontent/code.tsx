import React from 'react'

import { Props } from '../component_interface'

import SyntaxHighlighter from 'react-syntax-highlighter'

import { prism, tomorrow } from 'react-syntax-highlighter/dist/esm/styles/prism'

export function TCode({ node, theme }: Props) {
  return (
    <SyntaxHighlighter
      language={node.props.lang}
      style={theme === 'dark' ? tomorrow : prism}>
      {node.props.code}
    </SyntaxHighlighter>
  )
}
