import React, { Component } from 'react';

import './App.css';

import { updater } from './updater.js'
import { TComponent } from './components/factory.js';
import { Node } from './Nodes.js';
import { sessionValues } from './components/session.js';

function faviconTemplate(icon) {
  return `
    <svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 100 100%22>
      <text y=%22.9em%22 font-size=%2290%22>
        ${icon}
      </text>
    </svg>
  `.trim();
}

class App extends Component {
  rootNode() {
    const ret = new Node({})
    ret.props.name = 'container_component'
    ret.props.id = 'container_root'
    return ret
  }

  constructor(props) {
    super(props);
    this.state = {
      data: {
        page_names: [],
        page_confs: {},
      },
      nodes: {
        container_root: this.rootNode(),
      },
    }
    this.getPageData()
  }

  componentDidMount() {
    this.update({})
  }

  update(event) {
    updater(event, () => {
      this.setState({
        nodes: {
          container_root: this.rootNode(),
        },
      })
    }, () => {
      sessionValues = {}
    }, (data) => {
      const compID = data.component.id
      this.setState((prevState) => {
        const newNodes = { ...prevState.nodes }
        newNodes[data.container_id].children.push(compID)
        newNodes[compID] = new Node(data.component)
        return {
          nodes: newNodes,
        }
      })
    })
  }

  getPageData() {
    fetch('/api/pages').then(function (resp) {
      return resp.json()
    }).then(data => {
      const pageName = window.location.pathname.substring(1)
      const curconf = data.page_confs[pageName]
      document.title = curconf.title
      if (curconf.emoji) {
        const iconEle = document.querySelector(`head > link[rel='icon']`)
        iconEle.setAttribute(`href`,
          `data:image/svg+xml,${faviconTemplate(curconf.emoji)}`)
      }
      this.setState({ data })
    }).catch(e => {
      console.log(e)
    })
  }

  render() {
    return (
      <div>
        <nav class="navbar" role="navigation" aria-label="main navigation">
          <div class="navbar-menu container">
            <div class="navbar-start">
              {
                this.state.data.page_names.map(name =>
                  <a class="navbar-item" href={'/' + name}>
                    {this.state.data.page_confs[name].emoji}
                    {this.state.data.page_confs[name].title}
                  </a>
                )
              }
            </div>
            <div class="navbar-end">
              <div class="buttons">
                <a class="button navbar-item" onClick={() => { this.update({}) }}>
                  Rerun
                </a>
              </div>
            </div>
          </div>
        </nav>
        <div class="container">
          <TComponent node={this.state.nodes.container_root}
            update={(e) => { this.update(e) }}
            nodes={this.state.nodes} />
        </div>
      </div>
    )
  }
}

export default App;
