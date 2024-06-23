package main

import (
	"log"

	"github.com/mudream4869/toolgui/toolgui/component"
	"github.com/mudream4869/toolgui/toolgui/executor"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

func Main(s *framework.Session, c *framework.Container) error {
	if component.Button(s, c, "say") {
		component.Text(c, "hi")
	}

	if component.Button(s, c, "say2") {
		component.Text(c, "hi2")
	}
	return nil
}

func main() {
	e := executor.NewWebExecutor()
	e.AddPage("index", "Index", Main)
	log.Println("Starting service...")
	e.StartService(":3000")
}
