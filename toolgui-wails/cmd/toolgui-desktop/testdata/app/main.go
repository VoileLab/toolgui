//go:build externalfixture && cli_fixture

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"

	tgwails "github.com/voilelab/toolgui/toolgui-wails"
	"github.com/voilelab/toolgui/toolgui/tgcomp"
	"github.com/voilelab/toolgui/toolgui/tgframe"
)

var version = "unset"

func main() {
	if os.Getenv("TOOLGUI_DESKTOP_PROBE") == "1" {
		if len(os.Args) > 1 && os.Args[1] == "exit" {
			os.Exit(7)
		}
		index, err := tgwails.Assets.ReadFile("frontend/dist/index.html")
		if err != nil || len(index) == 0 || strings.Contains(string(index), "toolgui: assets not built") {
			log.Fatalf("invalid embedded frontend: %v", err)
		}
		info, _ := debug.ReadBuildInfo()
		tags := ""
		for _, setting := range info.Settings {
			if setting.Key == "-tags" {
				tags = setting.Value
			}
		}
		exe, _ := os.Executable()
		cwd, _ := os.Getwd()
		if err := json.NewEncoder(os.Stdout).Encode(map[string]any{
			"version": version, "tags": tags, "args": os.Args[1:], "exe": exe, "cwd": cwd,
		}); err != nil {
			log.Fatal(err)
		}
		return
	}

	app := tgframe.NewApp()
	app.AddPage("index", "Index", func(p *tgframe.Params) error {
		name := tgcomp.Textbox(p.Sidebar, "What's your name?")
		tgcomp.Text(p.Main, "Hello "+name)
		if tgcomp.Button(p.Main, "Keep going") {
			tgcomp.Text(p.Main, "The button works")
		}
		// An optional GUI smoke run exits after the built frontend has loaded
		// and called Start through Wails. Normal runs stay interactive.
		if os.Getenv("TOOLGUI_DESKTOP_GUI_SMOKE") == "1" {
			fmt.Println("desktop frontend started")
			time.AfterFunc(time.Second, func() { os.Exit(0) })
		}
		return nil
	})
	if err := tgwails.NewExecutor(app, &tgwails.Conf{Title: "ToolGUI Desktop CLI"}).Run(); err != nil {
		log.Fatal(err)
	}
}
