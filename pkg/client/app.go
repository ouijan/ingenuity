package client

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ark "github.com/mlange-42/ark/ecs"

	"github.com/ouijan/ingenuity/pkg/client/renderer"
	"github.com/ouijan/ingenuity/pkg/client/resources"
	"github.com/ouijan/ingenuity/pkg/client/systems"
	"github.com/ouijan/ingenuity/pkg/core/config"
	"github.com/ouijan/ingenuity/pkg/core/ecs"
	"github.com/ouijan/ingenuity/pkg/core/ecs/components"
	"github.com/ouijan/ingenuity/pkg/core/utils"
)

type ClientApp struct {
	exit   bool
	window *Window
	config *config.Config
	ecs    *ecs.EntityAdmin
	camera rl.Camera2D
}

func (a *ClientApp) Run() error {
	screenW := 1600 // 800
	screenH := 900  // 450
	a.window.Open(screenW, screenH, "Ingenuity Client")
	defer a.window.Close()

	a.SetupWorld()
	a.camera.Offset = rl.Vector2{
		X: float32(screenW) / 2,
		Y: float32(screenH) / 2,
	}
	a.camera.Zoom = .25

	a.ecs.Activate()
	for !a.exit && !rl.WindowShouldClose() {
		a.ecs.Update(rl.GetFrameTime())
		renderer.Render(a.camera)
		renderer.Clear()
	}

	return nil
}

func (a *ClientApp) SetupWorld() error {
	ark.AddResource(&a.ecs.World, a)
	ark.AddResource(&a.ecs.World, a.config)
	ark.AddResource(&a.ecs.World, a.window)
	ark.AddResource(&a.ecs.World, &a.camera)
	ark.AddResource(&a.ecs.World, &resources.UserInputStore{})

	ecs.AddSystem(a.ecs, &systems.TinkerSystem{})
	ecs.AddSystem(a.ecs, &systems.InputSystem{})
	ecs.AddSystem(a.ecs, &systems.InputHandlerSystem{})

	textFactory := ark.NewMap3[components.Metadata, components.Transform2D, components.Text](
		&a.ecs.World,
	)
	textFactory.NewEntity(
		&components.Metadata{Name: utils.MessageDisplayName},
		&components.Transform2D{X: 190, Y: 200},
		&components.Text{Content: "Message Display", FontSize: 20},
	)
	textFactory.NewEntity(
		&components.Metadata{Name: utils.FPSDisplayName},
		&components.Transform2D{X: 10, Y: 10},
		&components.Text{Content: "FPS Display", FontSize: 20},
	)

	netTextFactory := ark.NewMap4[components.Metadata, components.Transform2D, components.Text, components.NetworkedEntity](
		&a.ecs.World,
	)
	netTextFactory.NewEntity(
		&components.Metadata{Name: utils.ServerTickDisplayName},
		&components.Transform2D{X: 10, Y: 420},
		&components.Text{Content: "", FontSize: 20},
		components.NewNetworkedEntity(1),
	)
	return nil
}

func (a *ClientApp) Close() {
	a.exit = true
}

func NewClientApp(cfg *config.Config) *ClientApp {
	return &ClientApp{
		exit:   false,
		config: cfg,
		window: NewWindow(),
		ecs:    ecs.NewEntityAdmin(),
		camera: rl.Camera2D{
			Zoom: 1.0,
		},
	}
}
