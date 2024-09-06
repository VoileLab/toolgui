package main

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/mudream4869/toolgui/toolgui/tgcomp"
	"github.com/mudream4869/toolgui/toolgui/tgexec"
	"github.com/mudream4869/toolgui/toolgui/tgframe"
)

func Main(p *tgframe.Params) error {
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
	app := tgframe.NewApp()
	app.AddPage("main", "Main", Main)

	e := tgexec.NewWebExecutor(app)
	log.Println("Starting service...")
	e.StartService(":3000")
}
