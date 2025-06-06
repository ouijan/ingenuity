package server

import (
	"time"

	ark "github.com/mlange-42/ark/ecs"

	"github.com/ouijan/ingenuity/pkg/core/config"
	"github.com/ouijan/ingenuity/pkg/core/ecs"
	"github.com/ouijan/ingenuity/pkg/core/ecs/components"
	"github.com/ouijan/ingenuity/pkg/core/log"
	"github.com/ouijan/ingenuity/pkg/core/net"
	"github.com/ouijan/ingenuity/pkg/core/utils"
	"github.com/ouijan/ingenuity/pkg/server/systems"
)

const ServerTickRate = 160 * time.Millisecond

type ServerApp struct {
	config *config.Config
	server *net.Server
	ecs    *ecs.EntityAdmin
}

func (a *ServerApp) Init() error {
	err := a.server.Start()
	if err != nil {
		return err
	}

	ark.AddResource(&a.ecs.World, a.server)
	ark.AddResource(&a.ecs.World, a.config)

	ecs.AddSystem(a.ecs, &systems.ServerNetSync{})

	textFactory := ark.NewMap4[components.Metadata, components.Transform2D, components.Text, components.NetworkedEntity](
		&a.ecs.World,
	)
	textFactory.NewEntity(
		&components.Metadata{Name: utils.ServerTickDisplayName},
		&components.Transform2D{X: 10, Y: 420},
		&components.Text{Content: "Server Tick Here", FontSize: 20},
		components.NewNetworkedEntity(1),
	)
	return nil
}

func (a *ServerApp) Run() error {
	log.Info("Starting server")

	go a.server.Listen()
	defer a.server.Close()

	a.ecs.Activate()
	go a.gameLoop()
	go a.handleIncomingMessages()

	a.handleIncomingCommands()
	return nil
}

func (a *ServerApp) gameLoop() {
	ticker := time.NewTicker(ServerTickRate)
	lastTick := time.Now()

	for {
		t := <-ticker.C
		dt := t.Sub(lastTick).Seconds()
		lastTick = t
		a.ecs.Update(float32(dt))
	}
}

func (a *ServerApp) handleIncomingMessages() {
	for msg := range a.server.Buffer.SyncCh {
		log.Info("<- Received message: %s", msg)
	}
}

func (a *ServerApp) handleIncomingCommands() {
	for cmd := range utils.ReadStdIn {
		log.Info("Command >>> %s", cmd)
		if cmd == "exit" {
			return
		}

		log.Info("-> Sending message: %s", cmd)

		// msg := NetMsg{
		// 	Debug: &NetMsg_Debug{Msg: cmd},
		// }
		//
		// jd, err := json.Marshal(msg)
		// if err != nil {
		// 	log.Error("Failed to marshal message: %v", err)
		// 	return
		// }
		// a.server.Broadcast(jd)
	}
}

func NewServerApp(cfg *config.Config) *ServerApp {
	return &ServerApp{
		config: cfg,
		server: net.NewServer(cfg.Port),
		ecs:    ecs.NewEntityAdmin(),
	}
}
