package main

import (
	"log"
	"time"

	"github.com/mudream4869/toolgui/toolgui/component"
	"github.com/mudream4869/toolgui/toolgui/executor"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

func Main(s *framework.Session, c *framework.Container) error {
	pb := component.ProgressBar(c, 0, "Running...")
	for i := range 10 {
		pb.SetValue(i*10 + 10)
		time.Sleep(100 * time.Millisecond)
	}
	pb.SetLabel("OK.")
	time.Sleep(5 * time.Second)
	pb.Remove()
	return nil
}

func main() {
	e := executor.NewWebExecutor()
	e.AddPage("index", "Index", Main)
	log.Println("Starting service...")
	e.StartService(":3000")
}
