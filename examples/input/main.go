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

	val := 0

	if component.Checkbox(s, c, "a + b -> a - b") {
		val = a - b
	} else {
		val = a + b
	}

	selected := component.Select(s, c, "rate", []string{
		"x1",
		"x2",
		"x3",
	})

	switch selected {
	case "x1":
	case "x2":
		val *= 2
	case "x3":
		val *= 3
	default:
	}

	component.Text(c, fmt.Sprintf("Value = %d", val))

	return nil
}

func main() {
	e := executor.NewWebExecutor()
	e.AddPage("index", "Index", Main)
	log.Println("Starting service...")
	e.StartService()
}
