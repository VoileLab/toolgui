import Markdown from 'react-markdown'

export function TText({ node }) {
  return (
    <div>{node.props.text}</div>
  )
}

export function TDivider({ }) {
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

export function TMarkdown({ node }) {
  return (
    <div class="content">
      <Markdown>{node.props.text}</Markdown>
    </div>
  )
}