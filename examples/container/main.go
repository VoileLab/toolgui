package main

import (
	"log"
	"time"

	"github.com/mudream4869/toolgui/toolgui/component"
	"github.com/mudream4869/toolgui/toolgui/executor"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

func Main(s *framework.Session, c *framework.Container) error {
	startCont := c.AddContainer("start")
	component.Text(c, "1")
	time.Sleep(time.Second)

	component.Text(startCont, "2")

	return nil
}

func main() {
	e := executor.NewWebExecutor()
	e.AddPage("index", "Index", Main)
	log.Println("Starting service...")
	e.StartService()
}
