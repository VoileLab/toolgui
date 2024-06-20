import Markdown from 'react-markdown'
import hljs from 'highlight.js'
import 'highlight.js/styles/default.css'

export function TText({ node }) {
  return (
    <div>{node.props.text}</div>
  )
}

export function TDivider() {
  return (
    <hr />
  )
}

export function TTitle({ node }) {
  return (
    <h1 id={node.props.id} class="title">
      {node.props.text}
    </h1>
  )
}

export function TSubtitle({ node }) {
  return (
    <h2 id={node.props.id} class="subtitle">
      {node.props.text}
    </h2>
  )
}

export function TImage({ node }) {
  return (
    <img src={node.props.src} />
  )
}

export function TMarkdown({ node }) {
  return (
    <div class="content">
      <Markdown>{node.props.text}</Markdown>
    </div>
  )
}

export function TCode({ node }) {
  const highlightHTML = hljs.highlight(
    node.props.code,
    { language: node.props.lang }
  ).value
  return (
    <div class="content">
      <pre>
        <div dangerouslySetInnerHTML={{ __html: highlightHTML }} />
      </pre>
    </div>
  )
}
