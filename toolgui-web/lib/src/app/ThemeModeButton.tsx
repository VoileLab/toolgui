import React, { Component } from 'react'

interface ThemeModeButtonState {
  dark_mode: boolean
}

export class ThemeModeButton extends Component<{}, ThemeModeButtonState> {
  constructor(props: any) {
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
      <button className="button" onClick={() => { this.toggleTheme() }}>
        <span className="icon">
          {this.state.dark_mode ? <i className="fas fa-moon"></i> : <i className="fas fa-sun"></i>}
        </span>
      </button>
    )
  }
}