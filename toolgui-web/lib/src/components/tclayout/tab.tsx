import React, { useState } from 'react'
import { Props } from '../component_interface'
import { TComponent } from '../factory'

export function TTab({ node, update, upload }: Props) {
  const [activeTab, setActiveTab] = useState(node.props.tabs[0])

  const activeIndex = node.props.tabs.indexOf(activeTab)

  return (
    <>
      <div className="tabs">
        <ul>
          {
            node.props.tabs.map((tab: string) => (
              <li className={`${activeTab === tab ? 'is-active' : ''} is-boxed`}>
                <a onClick={() => setActiveTab(tab)}>{tab}</a>
              </li>
            ))
          }
        </ul>
      </div>
      <div>
        <TComponent node={node.children[activeIndex]}
          update={update}
          upload={upload} />
      </div>
    </>
  )
}
