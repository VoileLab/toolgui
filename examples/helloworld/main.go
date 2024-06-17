package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mudream4869/toolgui/toolgui/component"
	"github.com/mudream4869/toolgui/toolgui/executor"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

type Foo struct {
	Str  string `json:"str"`
	Int  int    `json:"int"`
	Bool bool   `json:"bool"`
	Null any    `json:"null"`
}

func Main(r *framework.Session, c *framework.Container) error {
	component.Title(c, "Hello world")
	component.Subtitle(c, "This is a hello-world example.")

	component.Divider(c)

	col1, col2 := component.Column2(c, "coltest")

	box := component.Box(col2, "box1")
	component.Text(box, "Box test")
	component.JSON(box, box)
	component.JSON(box, &Foo{
		Str:  "123",
		Int:  123,
		Bool: true,
		Null: nil,
	})

	component.Text(col1, "Please wait for 5 seconds...")

	for i := range 5 {
		time.Sleep(time.Second)
		component.Text(col1, fmt.Sprintf("%d second(s)...", i+1))
	}

	component.Info(col1, "Info", "ok.")

	return nil
}

func main() {
	e := executor.NewWebExecutor()
	e.AddPage("index", "Index", Main)
	log.Println("Starting service...")
	e.StartService()
}
