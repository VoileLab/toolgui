import React from 'react'

import hljs from 'highlight.js'
import 'highlight.js/styles/default.css'
import Markdown from 'react-markdown'
import { Props } from '../component_interface'

export function TMarkdown({ node }: Props) {
  return (
    <div className="content">
      <Markdown children={node.props.text}
        components={{
          code(props) {
            const { children, className, node, ...rest } = props
            const match = /language-(\w+)/.exec(className || '')
            if (!match) {
              return (
                <code {...rest} className={className}>
                  {children}
                </code>
              )
            }

            const lang = match[1]
            const highlightHTML = hljs.highlight(
              String(children).replace(/\n$/, ''),
              { language: lang }
            ).value

            return (
              <div dangerouslySetInnerHTML={{ __html: highlightHTML }} />
            )
          }
        }}
      />
    </div>
  )
}