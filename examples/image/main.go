package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/png"

	"log"

	"github.com/mudream4869/toolgui/toolgui/component"
	"github.com/mudream4869/toolgui/toolgui/executor"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

//go:embed logo192.png
var logo []byte

var pngLogo image.Image

func Main(s *framework.Session, c *framework.Container, _ *framework.Container) error {
	component.Title(c, "By image")
	component.Image(c, pngLogo)
	component.Title(c, "By url")
	component.ImageByURL(c, "https://placehold.co/100x100")
	return nil
}

func main() {
	var err error
	pngLogo, err = png.Decode(bytes.NewReader(logo))
	if err != nil {
		log.Panic(err)
	}

	e := executor.NewWebExecutor()
	e.AddPage("index", "Index", Main)
	log.Println("Starting service...")
	e.StartService(":3000")
}
