# ToolGUI

This Go package provides a framework for rapidly building interactive data
dashboards and web applications. It aims to offer a similar development
experience to Streamlit for Python users.

> [!WARNING]
> ⚠️ Under Development:
> 
> The API for this package is still under development,
> and may be subject to changes in the future.

## Hello world

```go
package main

import (
	"log"
	"github.com/mudream4869/toolgui/toolgui/component/tccontent"
	"github.com/mudream4869/toolgui/toolgui/component/tcinput"
	"github.com/mudream4869/toolgui/toolgui/executor"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

func Main(s *framework.Session, c *framework.Container, sidebar *framework.Container) error {
	name := tcinput.Textbox(s, sidebar, "What's your name?")
	if name != "" {
		tccontent.Text(sidebar, "Hi "+name+"~")
	}

	tccontent.Text(c, "hello ")
	if tcinput.Button(s, c, "keep going") {
		tccontent.Text(c, "world")
	}
	return nil
}

func main() {
	e := executor.NewWebExecutor()
	e.AddPage("index", "Index", Main)
	log.Println("Starting service...")
	e.StartService(":3000")
}
```
