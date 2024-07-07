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
