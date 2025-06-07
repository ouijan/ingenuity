package client

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ark "github.com/mlange-42/ark/ecs"

	"github.com/ouijan/ingenuity/pkg/client/systems"
	"github.com/ouijan/ingenuity/pkg/core/config"
	"github.com/ouijan/ingenuity/pkg/core/ecs"
	"github.com/ouijan/ingenuity/pkg/core/ecs/components"
	"github.com/ouijan/ingenuity/pkg/core/log"
	"github.com/ouijan/ingenuity/pkg/core/utils"
)

type ClientApp struct {
	exit         bool
	window       *Window
	config       *config.Config
	ecs          *ecs.EntityAdmin
	textEntities *ark.Filter3[components.Metadata, components.Transform2D, components.Text]
}

func (a *ClientApp) Init() error {
	ark.AddResource(&a.ecs.World, a)
	ark.AddResource(&a.ecs.World, a.config)
	ark.AddResource(&a.ecs.World, a.window)

	ecs.AddSystem(a.ecs, &systems.TinkerSystem{})

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

func (a *ClientApp) Run() error {
	log.Info("Starting client")

	a.window.Open(800, 450, "Ingenuity Client")
	defer a.window.Close()

	a.ecs.Activate()
	for !a.exit && !rl.WindowShouldClose() {
		a.update()
		a.render()
	}

	return nil
}

func (a *ClientApp) update() {
	dt := rl.GetFrameTime()
	a.ecs.Update(dt)
}

func (a *ClientApp) render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	query := a.textEntities.Query()

	for query.Next() {
		_, trans, text := query.Get()
		// log.Debug("Content: %s, Position: (%f, %f), FontSize: %d",
		// 	text.Content, trans.X, trans.Y, text.FontSize,
		// )
		rl.DrawText(
			text.Content,
			int32(trans.X),
			int32(trans.Y),
			int32(text.FontSize),
			rl.LightGray,
		)
	}

	rl.EndDrawing()
}

func (a *ClientApp) Close() {
	a.exit = true
}

func NewClientApp(cfg *config.Config) *ClientApp {
	ecs := ecs.NewEntityAdmin()
	textEntities := ark.NewFilter3[components.Metadata, components.Transform2D, components.Text](
		&ecs.World,
	)
	return &ClientApp{
		exit:         false,
		config:       cfg,
		window:       NewWindow(),
		ecs:          ecs,
		textEntities: textEntities,
	}
}
