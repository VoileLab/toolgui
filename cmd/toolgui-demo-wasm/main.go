//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"

	toolguidemo "github.com/mudream4869/toolgui/toolgui-demo"
)

func main() {
	app := toolguidemo.NewApp()
	js.Global().Set("getAppConf", js.FuncOf(func(this js.Value, args []js.Value) any {
		bs, err := json.Marshal(app.AppConf())
		if err != nil {
			panic(err)
		}
		return string(bs)
	}))

	select {}
}
