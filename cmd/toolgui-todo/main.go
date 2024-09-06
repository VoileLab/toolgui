package main

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/mudream4869/toolgui/toolgui/executor"
	"github.com/mudream4869/toolgui/toolgui/framework"
	"github.com/mudream4869/toolgui/toolgui/tgcomp/tccontent"
	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcinput"
)

func Main(p *framework.Params) error {
	tccontent.Title(p.Main, "Example for Todo App")

	var todos []string
	err := p.State.GetObject("todos", &todos)
	if err != nil {
		return err
	}

	inp := tcinput.Textbox(p.State, p.Main, "Add todo")
	if tcinput.Button(p.State, p.Main, "Add") {
		todos = append(todos, inp)
		p.State.Set("todos", todos)
	}

	for i, todo := range todos {
		tccontent.TextWithID(p.Main,
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
