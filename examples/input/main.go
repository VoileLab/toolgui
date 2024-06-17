package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/mudream4869/toolgui/toolgui/component"
	"github.com/mudream4869/toolgui/toolgui/executor"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

func Main(s *framework.Session, c *framework.Container) error {
	aStr := component.Textbox(s, c, "a")
	bStr := component.Textbox(s, c, "b")

	a, err := strconv.Atoi(aStr)
	if err != nil {
		return err
	}

	b, err := strconv.Atoi(bStr)
	if err != nil {
		return err
	}

	if component.Checkbox(s, c, "a + b -> a - b") {
		component.Text(c, fmt.Sprintf("a - b = %d", a-b))
	} else {
		component.Text(c, fmt.Sprintf("a + b = %d", a+b))
	}

	return nil
}

func main() {
	e := executor.NewWebExecutor()
	e.AddPage("index", "Index", Main)
	log.Println("Starting service...")
	e.StartService()
}
