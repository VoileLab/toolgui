# ToolGUI

This go package is trying to replicate the streamlit in golang.

Warning: This repo is in heavy developement. The API of the go package may change.

## Hello world

```go
package main

import (
	"log"
	"github.com/mudream4869/toolgui/toolgui/component"
	"github.com/mudream4869/toolgui/toolgui/executor"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

func Main(s *framework.Session, c *framework.Container) error {
    component.Text(c, "hello ")
    if component.Button(s, c, "keep going") {
        component.Text(c, "world")
    }
	return nil
}

func main() {
	e := executor.NewWebExecutor()
	e.AddPage("index", "Index", Main)
	log.Println("Starting service...")
	e.StartService()
}
```
