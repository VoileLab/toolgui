import React from 'react'
import { Props } from '../component_interface'

export function TImage({ node }: Props) {
  return (
    <img src={node.props.src} width={node.props.width} />
  )
}


