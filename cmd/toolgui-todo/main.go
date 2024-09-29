package main

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/VoileLab/toolgui/toolgui/tgcomp"
	"github.com/VoileLab/toolgui/toolgui/tgexec"
	"github.com/VoileLab/toolgui/toolgui/tgframe"
)

type TODOList struct {
	items []string
}

func (t *TODOList) add(item string) {
	t.items = append(t.items, item)
}

func Main(p *tgframe.Params) error {
	tgcomp.Title(p.Main, "Example for Todo App")

	todoList := p.State.Default("todoList", &TODOList{}).(*TODOList)

	inp := tgcomp.Textbox(p.State, p.Main, "Add todo")
	if tgcomp.Button(p.State, p.Main, "Add") {
		todoList.add(inp)
	}

	for i, todo := range todoList.items {
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
