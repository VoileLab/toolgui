# ToolGUI

[![Go Reference](https://pkg.go.dev/badge/github.com/mudream4869/toolgui.svg)](https://pkg.go.dev/github.com/mudream4869/toolgui)

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

func Main(p *framework.Params) error {
	name := tcinput.Textbox(p.State, p.Sidebar, "What's your name?")
	if name != "" {
		tccontent.Text(p.Sidebar, "Hi "+name+"~")
	}

	tccontent.Text(p.Main, "hello ")
	if tcinput.Button(p.State, p.Main, "keep going") {
		tccontent.Text(p.Main, "world")
	}
	return nil
}

func main() {
	app := framework.NewApp()
	app.AddPage("index", "Index", Main)
	e := executor.NewWebExecutor(app)
	log.Println("Starting service...")
	e.StartService(":3000")
}
```

## For Dev

### Dependency

* [yarn](https://yarnpkg.com/): Frontend
* [cypress](https://www.cypress.io/): E2E Testing
* [taskfile](https://taskfile.dev/): Task runner

### Run demo

```shell
task run_demo
```

### Run E2E Test

```shell
task run_demo
```

```shell
cd toolgui-e2e
cypress e2e:chrome
cypress e2e:firefox
```
