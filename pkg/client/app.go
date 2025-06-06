package client

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	ark "github.com/mlange-42/ark/ecs"

	"github.com/ouijan/ingenuity/pkg/core/config"
	"github.com/ouijan/ingenuity/pkg/core/ecs"
	"github.com/ouijan/ingenuity/pkg/core/ecs/components"
	"github.com/ouijan/ingenuity/pkg/core/ecs/systems"
	"github.com/ouijan/ingenuity/pkg/core/log"
	"github.com/ouijan/ingenuity/pkg/core/net"
	"github.com/ouijan/ingenuity/pkg/core/utils"
)

type ClientApp struct {
	exit         bool
	window       *Window
	config       *config.Config
	client       *net.Client
	ecs          *ecs.EntityAdmin
	textEntities *ark.Filter3[components.Metadata, components.Transform2D, components.Text]
}

func (a *ClientApp) Init() error {
	err := a.client.Connect()
	if err != nil {
		return err
	}

	ark.AddResource(&a.ecs.World, a)
	ark.AddResource(&a.ecs.World, a.client)
	ark.AddResource(&a.ecs.World, a.config)
	ark.AddResource(&a.ecs.World, a.window)

	ecs.AddSystem(a.ecs, &TestClientSystem{})
	ecs.AddSystem(a.ecs, &systems.ClientNetSync{})

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

	go a.client.Listen()
	defer a.client.Close()

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
	addr := fmt.Sprintf("%s:%d", cfg.Client.Host, cfg.Port)

	ecs := ecs.NewEntityAdmin()
	textEntities := ark.NewFilter3[components.Metadata, components.Transform2D, components.Text](
		&ecs.World,
	)
	return &ClientApp{
		exit:         false,
		config:       cfg,
		window:       NewWindow(),
		client:       net.NewClient(addr),
		ecs:          ecs,
		textEntities: textEntities,
	}
}

// ------ TEST SYSTEM ------

type TestClientSystem struct {
	ca            ark.Resource[ClientApp]
	textEntities  *ark.Filter2[components.Metadata, components.Text]
	netEntities   *ark.Filter1[components.NetworkedEntity]
	netComponents *ark.Map3[components.Metadata, components.Transform2D, components.Text]
}

func (s *TestClientSystem) OnCreate(ea *ecs.EntityAdmin) {
	s.ca = ark.NewResource[ClientApp](&ea.World)
	s.textEntities = ark.NewFilter2[components.Metadata, components.Text](&ea.World)
	s.netEntities = ark.NewFilter1[components.NetworkedEntity](&ea.World)
	s.netComponents = ark.NewMap3[components.Metadata, components.Transform2D, components.Text](
		&ea.World,
	)
}

func (s *TestClientSystem) Update(dt float32) error {
	ca := s.ca.Get()
	if ca == nil {
		return nil
	}

	fpsMsg := fmt.Sprintf("FPS: %d, DT: %f", rl.GetFPS(), rl.GetFrameTime())

	// TODO: this is only grabbing 1 message per frame

	// data := &NetMsg{}
	// msg, hasMsg := utils.ChanSelect(ca.client.MsgCh)
	// if hasMsg {
	// 	// log.Info("<- Received message %s", msg)
	// 	err := json.Unmarshal(msg.Payload, data)
	// 	if err != nil {
	// 		log.Error("Failed to unmarshal message: %v \n %v", err, msg.Payload)
	// 		return err
	// 	}
	// }

	query := s.textEntities.Query()
	for query.Next() {
		meta, text := query.Get()

		switch meta.Name {
		case utils.FPSDisplayName:
			text.Content = fpsMsg
			// case MessageDisplayName:
			// 	if hasMsg && data.Debug != nil {
			// 		text.Content = data.Debug.Msg
			// 	}
		}
	}

	// if hasMsg && data.EntityUpdate != nil {
	// 	query := s.netEntities.Query()
	// 	for query.Next() {
	// 		net := query.Get()
	// 		if net.Id == data.EntityUpdate.Id {
	// 			s.updateEntity(data.EntityUpdate, query.Entity())
	// 		}
	// 	}
	// }

	return nil
}

// func (s *TestClientSystem) updateEntity(update *NetMsg_EntityUpdate, entity ark.Entity) {
// 	meta, trans, text := s.netComponents.Get(entity)
// 	if update.Metadata != nil {
// 		meta.Name = update.Metadata.Name
// 		meta.Tags = update.Metadata.Tags
// 	}
// 	if update.Transform != nil {
// 		trans.X = update.Transform.X
// 		trans.Y = update.Transform.Y
// 	}
// 	if update.Text != nil {
// 		text.Content = update.Text.Content
// 		text.FontSize = update.Text.FontSize
// 		text.Colour = update.Text.Colour
// 	}
// }

func (s *TestClientSystem) OnDestroy() {
	// Cleanup if needed
}

var _ ecs.System = &TestClientSystem{}
