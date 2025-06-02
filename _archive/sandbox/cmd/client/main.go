package main

import (
	"encoding/json"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/ouijan/ingenuity/pkg/app"
	"github.com/ouijan/ingenuity/pkg/core"
	"github.com/ouijan/ingenuity/pkg/engine"
	"github.com/ouijan/ingenuity/pkg/net"
	"github.com/ouijan/ingenuity/pkg/renderer"
	"github.com/ouijan/ingenuity/sandbox/src/pong"
)

func main() {
	c := app.Config{
		Name: "Pong Client",
	}
	a := app.NewApp(c)
	defer a.Close()

	a.Attach(&ClientNetworkLayer{})
	a.Attach(&ClientPresentationLayer{})

	if err := a.Run(30); err != nil {
		fmt.Println(err)
	}
}

// --------- Client Network Layer ---------

type ClientNetworkLayer struct {
	client net.Client
}

// OnAttach implements app.AppLayer.
func (c *ClientNetworkLayer) OnAttach() {
	c.client = net.NewClient("localhost:4223")
	err := c.client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	go c.client.Listen()
	go clientListen(c.client)

	core.OnEvent("client.input", func(e core.Event[engine.InputActionsEvent]) error {
		data, err := json.Marshal(e.Data.Actions)
		if err != nil {

			fmt.Println(err)
			return err
		}
		err = c.client.Write(data)
		if err != nil {
			fmt.Println(err)
		}
		return err
	})

	// err = c.client.Write([]byte("Hello, World!"))
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

// OnDetach implements app.AppLayer.
func (c *ClientNetworkLayer) OnDetach() {
	defer c.client.Close()
}

// OnUpdate implements app.AppLayer.
func (c *ClientNetworkLayer) OnUpdate(dt float32) error {
	return nil
}

var _ app.Layer = &ClientNetworkLayer{}

func clientListen(client net.Client) {
	for msg := range client.Read() {
		fmt.Printf("[%s]: %s\n", msg.GetAddr(), msg.GetPayload())
	}
}

// Client Input Layer

type ClientPresentationLayer struct {
	imc *engine.InputMappingContext
}

// OnAttach implements app.Layer.
func (c *ClientPresentationLayer) OnAttach() {
	core.AttachEventLogger(true, false)
	core.AttachRaylibLogger(true)
	core.AttachTraceLogger(true)

	engine.Window.Open(800, 450, "Ingenuity")

	c.imc = pong.NewPlayerInputContext()
	engine.Input.Register(c.imc)

	// TODO: Setup action mappings so that we can check data being sent to the server
	core.OnEvent("client.input", func(e core.Event[engine.InputActionsEvent]) error {
		jsonData, err := json.Marshal(e.Data)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(string(jsonData))
		return nil
	})
}

// OnDetach implements app.Layer.
func (c *ClientPresentationLayer) OnDetach() {
	engine.Window.Close()
	engine.Input.Unregister(c.imc)
}

// OnUpdate implements app.Layer.
func (c *ClientPresentationLayer) OnUpdate(dt float32) error {
	if rl.WindowShouldClose() {
		core.EmitEventSync("app.exit", &core.Event[any]{})
		return nil
	}

	renderer.AddCall(5, 0, func() {
		rl.DrawFPS(10, 10)
	})
	renderer.Render(engine.Window.CanvasWidth, engine.Window.CanvasHeight)

	core.EmitEvent("client.input", engine.InputActionsEvent{
		Actions: engine.Input.GetAll(),
	})
	// engine.Input.Update()
	return nil
}

var _ app.Layer = &ClientPresentationLayer{}
