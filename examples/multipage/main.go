package main

import (
	"log"

	"github.com/mudream4869/toolgui/toolgui/component"
	"github.com/mudream4869/toolgui/toolgui/executor"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

func Page1(_ *framework.Session, c *framework.Container) error {
	c.AddComponent(component.NewTextComponent("page1"))
	return nil
}

func Page2(_ *framework.Session, c *framework.Container) error {
	c.AddComponent(component.NewTextComponent("page2"))
	return nil
}

func main() {
	e := executor.NewWebExecutor()
	e.AddPageByConfig(&executor.PageConfig{
		Name:  "page1",
		Title: "Page1",
		Emoji: "🐱",
	}, Page1)
	e.AddPageByConfig(&executor.PageConfig{
		Name:  "page2",
		Title: "Page2",
		Emoji: "🔄",
	}, Page2)
	log.Println("Starting service...")
	e.StartService(":3000")
}
