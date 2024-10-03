import React, { Component } from 'react'

interface ThemeModeButtonState {
  darkMode: string
}

interface ThemeModeButtonProps {
  onChange: (darkMode: string) => void
}

export class ThemeModeButton extends Component<ThemeModeButtonProps, ThemeModeButtonState> {
  constructor(props: ThemeModeButtonProps) {
    super(props);
    this.state = {
      darkMode: localStorage.darkMode,
    }
  }

  componentDidMount() {
    const root = document.getElementsByTagName('html')[0];
    root.className = 'theme-' + this.state.darkMode
  }

  toggleTheme() {
    this.setState((preState) => {
      const newValue = preState.darkMode === 'dark' ? 'light' : 'dark'
      const root = document.getElementsByTagName('html')[0];
      root.className = 'theme-' + newValue
      localStorage.darkMode = newValue
      this.props.onChange(newValue)
      return {
        darkMode: newValue
      }
    })
  }

  render() {
    return (
      <button className="button" onClick={() => { this.toggleTheme() }}>
        <span className="icon">
          {this.state.darkMode === 'dark' ?
            <i className="fas fa-moon"></i> :
            <i className="fas fa-sun"></i>}
        </span>
      </button>
    )
  }
}