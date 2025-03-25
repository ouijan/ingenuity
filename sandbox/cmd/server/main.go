package main

import (
	"fmt"

	"github.com/ouijan/ingenuity/pkg/app"
	"github.com/ouijan/ingenuity/pkg/core"
	"github.com/ouijan/ingenuity/pkg/engine"
	"github.com/ouijan/ingenuity/pkg/net"

	"github.com/ouijan/ingenuity/sandbox/src/pong"
)

func main() {
	c := app.Config{
		Name: "Pong Server",
	}
	a := app.NewApp(c)
	a.Attach(&ServerNetworkLayer{})
	a.Attach(&WorldLayer{})
	defer a.Close()

	if err := a.Run(30); err != nil {
		fmt.Println(err)
	}
}

// --------- Network Layer ---------

type ServerNetworkLayer struct {
	server net.Server
}

// OnAttach implements app.AppLayer.
func (n *ServerNetworkLayer) OnAttach() {
	n.server = net.NewServer("localhost:4223", 2)
	go n.handleUserConnected()
	go n.handleUserDisconnected()
	go n.readMessages()
	go n.startServer()
}

// OnDetach implements app.AppLayer.
func (n *ServerNetworkLayer) OnDetach() {
	defer n.server.Close()
}

// OnUpdate implements app.AppLayer.
func (n *ServerNetworkLayer) OnUpdate(dt float32) error {
	return nil
}

func (n *ServerNetworkLayer) startServer() {
	err := n.server.Listen()
	if err != nil {
		fmt.Println(err)
	}
}

func (n *ServerNetworkLayer) readMessages() {
	for msg := range n.server.Read() {
		core.EmitEvent("engine.network.message", msg)
	}
}

func (n *ServerNetworkLayer) handleUserConnected() {
	for addr := range n.server.OnConnect() {
		core.EmitEvent("engine.network.userConnected", addr)
	}
}

func (n *ServerNetworkLayer) handleUserDisconnected() {
	for addr := range n.server.OnDisconnect() {
		core.EmitEvent("engine.network.userDisconnected", addr)
	}
}

var _ app.Layer = &ServerNetworkLayer{}

// --------- World Layer ---------

type WorldLayer struct {
	world   *engine.World
	systems *engine.SystemManager
	netCh   chan net.Message
}

// OnAttach implements app.AppLayer.
func (w *WorldLayer) OnAttach() {
	w.netCh = make(chan net.Message, 100)
	core.OnEventCh("engine.network.message", w.netCh)

	w.world = engine.NewWorld()
	setupWorld(w.world)

	w.systems = engine.NewSystemManager()
	w.systems.Register(pong.NewPlayerSpawnSystem())
	w.systems.Register(pong.NewPlayerController())
	w.systems.Register(engine.NewPhysics2DSystem())
	// sm.Register(pong.NewGameplaySystem())

	// Network Notify System
}

// OnDetach implements app.AppLayer.
func (w *WorldLayer) OnDetach() {
	// panic("unimplemented")
}

// OnUpdate implements app.AppLayer.

func (w *WorldLayer) OnUpdate(dt float32) error {
	core.ReadCh(w.netCh, w.handleNetworkMessage)
	w.systems.Update(w.world, dt)
	return nil
}

func (w *WorldLayer) handleNetworkMessage(msg net.Message) error {
	// fmt.Printf("Processing -> [%s]: %s\n", msg.GetAddr(), msg.GetPayload())
	return nil
}

var _ app.Layer = &WorldLayer{}

func setupWorld(w *engine.World) {
	// Set up the world
	width := 800.0
	height := 450.0
	hWidth := width / 2
	hHeight := height / 2

	// Add Player 1
	p1 := engine.AddEntity(w)
	engine.AddComponent(w, p1, &pong.PlayerControllerComponent{
		OwnerId: "",
	})
	engine.AddComponent(w, p1, &pong.PaddleComponent{
		Speed: 2,
	})
	engine.AddComponent(w, p1, &engine.TransformComponent{
		X: 20, Y: hHeight,
	})
	engine.AddComponent(w, p1, &engine.BoxCollider2DComponent{
		T: 30, B: 30, L: 5, R: 5,
		Category:     1,
		CategoryMask: 1,
	})
	engine.AddComponent(w, p1, &engine.RigidBody2DComponent{
		Type: engine.RB_Kinematic,
		Mass: 1,
		Vx:   0,
		Vy:   0,
	})

	// Add Player 2
	p2 := engine.AddEntity(w)
	engine.AddComponent(w, p2, &pong.PlayerControllerComponent{
		OwnerId: "",
	})
	engine.AddComponent(w, p2, &pong.PaddleComponent{
		Speed: 2,
	})
	engine.AddComponent(w, p2, &engine.TransformComponent{
		X: float64(engine.Window.CanvasWidth - 20),
		Y: hHeight,
	})
	engine.AddComponent(w, p2, &engine.BoxCollider2DComponent{
		T: 30, B: 30, L: 5, R: 5,
		Category:     1,
		CategoryMask: 1,
	})
	engine.AddComponent(w, p2, &engine.RigidBody2DComponent{
		Type: engine.RB_Kinematic,
		Mass: 1,
		Vx:   0,
		Vy:   0,
	})

	// Add Ball
	ballSpeed := 2.0
	ball := engine.AddEntity(w)
	engine.AddComponent(w, ball, &pong.BallComponent{
		Speed: ballSpeed,
	})
	engine.AddComponent(
		w,
		ball,
		&engine.TransformComponent{X: hWidth, Y: hHeight},
	)
	engine.AddComponent(w, ball, &engine.BoxCollider2DComponent{
		T: 5, B: 5, L: 5, R: 5,
		Category:     1,
		CategoryMask: 1,
	})
	engine.AddComponent(w, ball, &engine.RigidBody2DComponent{
		Type: engine.RB_Dynamic,
		Mass: 1,
		// Vx:   -100 * ballSpeed,
		// Vy:   -100 * ballSpeed,
	})

	// Add Walls
	topWall := engine.AddEntity(w)
	engine.AddComponent(w, topWall, &engine.TransformComponent{
		X: hWidth,
		Y: 0,
	})
	engine.AddComponent(w, topWall, &engine.BoxCollider2DComponent{
		T: 5, B: 5, L: hWidth, R: hWidth,
		Category:     1,
		CategoryMask: 1,
	})
	engine.AddComponent(w, topWall, &engine.RigidBody2DComponent{
		Type: engine.RB_Static,
		Mass: 1,
		Vx:   0,
		Vy:   0,
	})

	bottomWall := engine.AddEntity(w)
	engine.AddComponent(w, bottomWall, &engine.TransformComponent{
		X: hWidth,
		Y: height,
	})
	engine.AddComponent(w, bottomWall, &engine.BoxCollider2DComponent{
		T: 5, B: 5, L: hWidth, R: hWidth,
		Category:     1,
		CategoryMask: 1,
	})
	engine.AddComponent(w, bottomWall, &engine.RigidBody2DComponent{
		Type: engine.RB_Static,
		Mass: 1,
		Vx:   0,
		Vy:   0,
	})

	// Add win triggers (left and right)
	leftWall := engine.AddEntity(w)
	engine.AddComponent(w, leftWall, &pong.GoalComponent{Id: pong.GoalLeft})
	engine.AddComponent(w, leftWall, &engine.TransformComponent{
		X: 0, Y: hHeight,
	})
	engine.AddComponent(w, leftWall, &engine.BoxCollider2DComponent{
		T: hHeight, B: hHeight, L: 5, R: 5,
		Category:     1,
		CategoryMask: 1,
	})

	rightWall := engine.AddEntity(w)
	engine.AddComponent(w, rightWall, &pong.GoalComponent{Id: pong.GoalRight})
	engine.AddComponent(w, rightWall, &engine.TransformComponent{
		X: width, Y: hHeight,
	})
	engine.AddComponent(w, rightWall, &engine.BoxCollider2DComponent{
		T: hHeight, B: hHeight, L: 5, R: 5,
		Category:     1,
		CategoryMask: 1,
	})
}
