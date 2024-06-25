import React, { Component } from 'react'

import './App.css';

import { updater } from './updater.js'
import { TComponent } from './components/factory.js'
import { Node } from './Nodes.js'
import { sessionValues } from './components/session.js'

const NOTIFY_TYPE_CREATE = 1
const NOTIFY_TYPE_UPDATE = 2
const NOTIFY_TYPE_DELETE = 3

function faviconTemplate(icon) {
  return `
    <svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 100 100%22>
      <text y=%22.9em%22 font-size=%2290%22>
        ${icon}
      </text>
    </svg>
  `.trim();
}

class ThemeModeButton extends Component {
  constructor(props) {
    super(props);
    this.state = {
      dark_mode: localStorage.darkMode === 'true'
    }
  }

  componentDidMount() {
    const root = document.getElementsByTagName('html')[0];
    root.className = this.state.dark_mode ? 'theme-dark' : 'theme-light';
  }

  toggleTheme() {
    this.setState((preState) => {
      const newValue = !preState.dark_mode
      const root = document.getElementsByTagName('html')[0];
      root.className = newValue ? 'theme-dark' : 'theme-light';
      localStorage.darkMode = newValue;
      return {
        dark_mode: newValue
      }
    })
  }

  render() {
    return (
      <button class="button" onClick={() => { this.toggleTheme() }}>
        <span class="icon">
          {this.state.dark_mode ? <i class="fas fa-moon"></i> : <i class="fas fa-sun"></i>}
        </span>
      </button>
    )
  }
}

class App extends Component {
  rootNode() {
    const ret = new Node({})
    ret.props.name = 'container_component'
    ret.props.id = 'container_component_container_root'
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
        container_component_container_root: this.rootNode(),
      },
      running: false,
      page_found: true,
      page_name: window.location.pathname.substring(1),
    }
    this.getPageData()
  }

  componentDidMount() {
    this.update({})
  }

  update(event) {
    updater(event, () => {
      this.setState((prevState) => {
        const newNodes = { ...prevState.nodes }
        for (const nodeID in newNodes) {
          if (nodeID == 'container_component_container_root') {
            continue
          }

          newNodes[nodeID].removing = true
        }

        return {
          running: true,
          nodes: newNodes,
        }
      })
    }, () => {
      sessionValues = {}
    }, (pack) => {
      switch (pack.type) {
        case NOTIFY_TYPE_CREATE: {
          const compID = pack.component.id

          if (compID in this.state.nodes && !this.state.nodes[compID].removing) {
            console.error('Depulicated component id:', compID)
            return
          }

          this.setState((prevState) => {
            const newNodes = { ...prevState.nodes }

            if (compID in this.state.nodes) {
              const parentID = this.state.nodes[compID].parentID
              const idx = newNodes[parentID].children.indexOf(compID)
              newNodes[parentID].children.splice(idx, 1)
            }

            const parentNode = newNodes[pack.container_id]
            var idx = 0
            for (var i = 0; i < parentNode.children.length; i++) {
              const prevCompID = parentNode.children[i]
              if (newNodes[prevCompID].removing) {
                break
              }
              idx = i + 1
            }
            parentNode.children.splice(idx, 0, compID)

            const oldNode = newNodes[compID]
            if (oldNode) {
              oldNode.props = pack.component
              oldNode.removing = false
            } else {
              newNodes[compID] = new Node(pack.component)
            }

            newNodes[compID].parentID = pack.container_id

            return {
              nodes: newNodes,
            }
          })
          break
        }
        case NOTIFY_TYPE_UPDATE: {
          const compID = pack.component.id
          if (!(compID in this.state.nodes)) {
            console.error('Try to update a node that doesn\'t exist:', compID)
            return
          }

          this.setState((prevState) => {
            const newNodes = { ...prevState.nodes }
            newNodes[compID].props = pack.component
            return {
              nodes: newNodes,
            }
          })
          break
        }
        case NOTIFY_TYPE_DELETE: {
          const compID = pack.component_id
          if (!(compID in this.state.nodes)) {
            console.error('Try to remove a node that doesn\'t exist:', compID)
            return
          }

          this.setState((prevState) => {
            const newNodes = { ...prevState.nodes }
            const parentID = this.state.nodes[compID].parentID
            const idx = newNodes[parentID].children.indexOf(compID)
            newNodes[parentID].children.splice(idx, 1)
            delete newNodes[compID]
            return {
              nodes: newNodes,
            }
          })
          break
        }
        default:
          console.error('Notify pack type error', pack.type)
      }
    }, (data) => {
      this.setState((prevState) => {
        const newNodes = {}
        for (const [key, node] of Object.entries(prevState.nodes)) {
          if (node.removing) {
            continue
          }

          newNodes[key] = node
        }

        for (const [key, node] of Object.entries(newNodes)) {
          node.children = node.children.filter((nodeID) => {
            return newNodes[nodeID] !== undefined
          })
        }

        return {
          running: false,
          nodes: newNodes,
        }
      })
    })
  }

  setIcon(emoji) {
    const iconEle = document.querySelector(`head > link[rel='icon']`)
    iconEle.setAttribute(`href`,
      `data:image/svg+xml,${faviconTemplate(emoji)}`)
  }

  getPageData() {
    fetch('/api/pages').then(resp => resp.json()).then(data => {
      const curconf = data.page_confs[this.state.page_name]
      if (curconf) {
        document.title = curconf.title
        if (curconf.emoji) {
          this.setIcon(curconf.emoji)
        }
      } else {
        document.title = 'Page not found'
        this.setIcon('❓')
        this.setState({ page_found: false })
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
                  <a className={`navbar-item ${name === this.state.page_name ? 'is-active' : ''}`} href={'/' + name}>
                    {this.state.data.page_confs[name].emoji}
                    {this.state.data.page_confs[name].title}
                  </a>
                )
              }
            </div>
            <div class="navbar-end">
              {this.state.running ?
                <div class="navbar-brand navbar-item">
                  <span class="icon">
                    <i class="fas fa-spinner fa-pulse"></i>
                  </span>
                </div> : ''}
              <div class="buttons">
                {this.state.page_found ?
                  <button class="button navbar-item" onClick={() => { this.update({}) }}>
                    Rerun
                  </button> : ''}
                <ThemeModeButton />
              </div>
            </div>
          </div>
        </nav>
        <div class="container">
          {this.state.page_found ?
            <TComponent node={this.state.nodes.container_component_container_root}
              update={(e) => { this.update(e) }}
              nodes={this.state.nodes} /> :
            <div class="columns is-centered">
              <div class="column is-three-quarters">
                <article class="message is-warning">
                  <div class="message-header">
                    <p>Oops! Page not found.</p>
                  </div>
                  <div class="message-body">
                    We're sorry, the page you requested was not found.
                    Try using the navigation menu to find what you're looking for.
                  </div>
                </article>
              </div>
            </div>
          }
        </div>
      </div>
    )
  }
}

export default App;
