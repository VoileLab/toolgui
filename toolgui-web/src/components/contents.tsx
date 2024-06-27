import Markdown from 'react-markdown'
import hljs from 'highlight.js'
import 'highlight.js/styles/default.css'
import { Props } from './component_interface'

export function TText({ node }: Props) {
  return (
    <div>{node.props.text}</div>
  )
}

export function TDivider() {
  return (
    <hr />
  )
}

export function TTitle({ node }: Props) {
  return (
    <h1 id={node.props.id} className="title">
      {node.props.text}
    </h1>
  )
}

export function TSubtitle({ node }: Props) {
  return (
    <h2 id={node.props.id} className="subtitle">
      {node.props.text}
    </h2>
  )
}

export function TImage({ node }: Props) {
  return (
    <img src={node.props.src} />
  )
}

export function TMarkdown({ node }: Props) {
  return (
    <div className="content">
      <Markdown>{node.props.text}</Markdown>
    </div>
  )
}

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
