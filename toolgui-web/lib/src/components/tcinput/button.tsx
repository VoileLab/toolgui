import React from "react"

import { Props } from "../component_interface"

export function TButton({ node, update }: Props) {
  var className = 'button'

  if (node.props.color) {
    className += ' is-' + node.props.color
  }

  return (
    <button
      id={node.props.id}
      className={className}
      disabled={node.props.disabled}
      onClick={
        (event: React.MouseEvent<HTMLButtonElement>) => {
          const target = event.target as HTMLButtonElement
          update({
            id: target.id,
            value: true,
            is_temp: true,
          })
        }
      }>
      {node.props.label}
    </button>
  )
}
