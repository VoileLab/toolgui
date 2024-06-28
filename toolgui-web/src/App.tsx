import React, { Component } from 'react'

import './App.css';

import { updater, initHealthSock } from './updater'
import { TComponent } from './components/factory'
import { ThemeModeButton } from './ThemeModeButton';
import { Forest } from './Nodes'
import { Event } from './components/component_interface';
import { clearSession } from './components/session'
import { MessagePageNotFound } from './MessagePageNotFound';

const NOTIFY_TYPE_CREATE = 1
const NOTIFY_TYPE_UPDATE = 2
const NOTIFY_TYPE_DELETE = 3

function faviconTemplate(icon: string) {
  return `
    <svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 100 100%22>
      <text y=%22.9em%22 font-size=%2290%22>
        ${icon}
      </text>
    </svg>
  `.trim();
}

interface AppState {
  data: {
    page_names: string[]
    page_confs: { [page_name: string]: any }
  }
  forest: Forest,
  running: boolean
  page_found: boolean
  page_name: string
}

class App extends Component<{}, AppState> {
  constructor(props: any) {
    super(props);
    this.state = {
      data: {
        page_names: [],
        page_confs: {},
      },
      forest: new Forest([
        'container_component_container_sidebar',
        'container_component_container_root',
      ]),
      running: false,
      page_found: true,
      page_name: window.location.pathname.substring(1),
    }
    this.getPageData()
  }

  componentDidMount() {
    this.update({})
    initHealthSock()
  }

  startUpdating() {
    this.setState((prevState) => {
      const newForest = prevState.forest.swallowCopy()
      newForest.setToRemoving()

      return {
        running: true,
        forest: newForest,
      }
    })
  }

  receiveNotifyPack(pack: any) {
    switch (pack.type) {
      case NOTIFY_TYPE_CREATE: {
        this.setState((prevState) => {
          const newForest = prevState.forest.swallowCopy()
          newForest.createNode(pack.component, pack.container_id)

          return {
            forest: newForest,
          }
        })
        break
      }
      case NOTIFY_TYPE_UPDATE: {
        this.setState((prevState) => {
          const newForest = prevState.forest.swallowCopy()
          newForest.updateNode(pack.component)
          return {
            forest: newForest,
          }
        })
        break
      }
      case NOTIFY_TYPE_DELETE: {
        this.setState((prevState) => {
          const newForest = prevState.forest.swallowCopy()
          newForest.removeNode(pack.component_id)

          return {
            forest: newForest,
          }
        })
        break
      }
      default:
        console.error('Notify pack type error', pack.type)
    }
  }

  finishUpdating() {
    this.setState((prevState) => {
      const newForest = prevState.forest.swallowCopy()
      newForest.removeNodeWithRemovingTag()

      return {
        running: false,
        forest: newForest,
      }
    })
  }

  update(event: Event) {
    updater(event,
      () => { this.startUpdating() },
      clearSession,
      (pack) => { this.receiveNotifyPack(pack) },
      () => { this.finishUpdating() })
  }

  setIcon(emoji: string) {
    const iconEle = document.querySelector(`head > link[rel='icon']`) as Element
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
        <nav className="navbar" role="navigation" aria-label="main navigation"
          style={{ position: 'fixed', top: 0, width: '100%' }}>
          <div className="navbar-menu container">
            <div className="navbar-start">
              {
                this.state.data.page_names.map(name =>
                  <a className={`navbar-item ${name === this.state.page_name ? 'is-active' : ''}`} href={'/' + name}>
                    {this.state.data.page_confs[name].emoji}
                    {this.state.data.page_confs[name].title}
                  </a>
                )
              }
            </div>
            <div className="navbar-end">
              {this.state.running ?
                <div className="navbar-brand navbar-item">
                  <span className="icon">
                    <i className="fas fa-spinner fa-pulse"></i>
                  </span>
                </div> : ''}
              <div className="buttons">
                {this.state.page_found ?
                  <button className="button navbar-item" onClick={() => { this.update({}) }}>
                    Rerun
                  </button> : ''}
                <ThemeModeButton />
              </div>
            </div>
          </div>
        </nav>
        <div className="container" style={{ 'paddingTop': '60px' }}>
          {this.state.page_found ?
            <section className="columns is-fullheight">
              {this.state.forest.nodes.container_component_container_sidebar.children.length > 0 ?
                <aside className="column is-2">
                  <div style={{ position: 'sticky', overflow: 'auto', top: '60px' }}>
                    <TComponent node={this.state.forest.nodes.container_component_container_sidebar}
                      update={(e) => { this.update(e) }}
                      nodes={this.state.forest.nodes} />
                  </div>
                </aside> : ''}
              <div className="column">
                <TComponent node={this.state.forest.nodes.container_component_container_root}
                  update={(e) => { this.update(e) }}
                  nodes={this.state.forest.nodes} />
              </div>
            </section>
            : <MessagePageNotFound />}
        </div>
      </div >
    )
  }
}

export default App;
