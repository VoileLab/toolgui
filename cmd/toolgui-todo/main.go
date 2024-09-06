package main

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/mudream4869/toolgui/toolgui/executor"
	"github.com/mudream4869/toolgui/toolgui/framework"
	"github.com/mudream4869/toolgui/toolgui/tgcomp"
)

func Main(p *framework.Params) error {
	tgcomp.Title(p.Main, "Example for Todo App")

	var todos []string
	err := p.State.GetObject("todos", &todos)
	if err != nil {
		return err
	}

	inp := tgcomp.Textbox(p.State, p.Main, "Add todo")
	if tgcomp.Button(p.State, p.Main, "Add") {
		todos = append(todos, inp)
		p.State.Set("todos", todos)
	}

	for i, todo := range todos {
		tgcomp.TextWithID(p.Main,
			fmt.Sprintf("%d: %s", i, todo),
			fmt.Sprintf("todo_%d", i))
	}

	return nil
}

func main() {
	app := framework.NewApp()
	app.AddPage("main", "Main", Main)

	e := executor.NewWebExecutor(app)
	log.Println("Starting service...")
	e.StartService(":3000")
}
