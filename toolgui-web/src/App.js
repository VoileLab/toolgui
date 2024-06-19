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
                  <a className={`navbar-item ${name == this.state.page_name ? 'is-active' : ''}`} href={'/' + name}>
                    {this.state.data.page_confs[name].emoji}
                    {this.state.data.page_confs[name].title}
                  </a>
                )
              }
            </div>
            <div class="navbar-end">
              <div class="buttons">
                {this.state.page_found ?
                  <a class="button navbar-item" onClick={() => { this.update({}) }}>
                    Rerun
                  </a> : ''}
              </div>
            </div>
          </div>
        </nav>
        <div class="container">
          {this.state.page_found ?
            <TComponent node={this.state.nodes.container_root}
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
