import Markdown from 'react-markdown'
import { Props } from '../component_interface'

export function TMarkdown({ node }: Props) {
  return (
    <div className="content">
      <Markdown>{node.props.text}</Markdown>
    </div>
  )
}