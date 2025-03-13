package pong

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ouijan/ingenuity/pkg/engine"
	"github.com/ouijan/ingenuity/pkg/renderer"
)

// ------------------ Scene ------------------

type pongScene struct {
	systems            []engine.System
	playerInputContext *engine.InputMappingContext
	playerController   *playerController
}

// Load implements engine.IScene.
func (p *pongScene) Load() {
	// panic("unimplemented")
}

// OnEnter implements engine.IScene.
func (p *pongScene) OnEnter(world *engine.World) {
	// Set Player Controller and Register Systems
	engine.Input.Register(p.playerInputContext)
	engine.Systems.Register(p.systems...)

	// Add Player
	player := engine.AddEntity(world)
	engine.AddComponent(world, player, &paddleComponent{})
	engine.AddComponent(world, player, &engine.TransformComponent{
		X: 20,
		Y: float64(engine.Window.CanvasHeight / 2),
	})
	engine.AddComponent(world, player, &engine.BoxCollider2DComponent{
		T: 30, B: 30, L: 5, R: 5,
		Category:     1,
		CategoryMask: 1,
	})
	engine.AddComponent(world, player, &engine.RigidBody2DComponent{
		Type: engine.RB_Kinematic,
		Mass: 1,
		Vx:   0,
		Vy:   0,
	})
	p.playerController.SetEntity(player)

	// Add Enemy
	enemy := engine.AddEntity(world)
	engine.AddComponent(world, enemy, &paddleComponent{})
	engine.AddComponent(world, enemy, &engine.TransformComponent{
		X: float64(engine.Window.CanvasWidth - 20),
		Y: float64(engine.Window.CanvasHeight / 2),
	})
	engine.AddComponent(world, enemy, &engine.BoxCollider2DComponent{
		T: 30, B: 30, L: 5, R: 5,
		Category:     1,
		CategoryMask: 1,
	})
	engine.AddComponent(world, enemy, &engine.RigidBody2DComponent{
		Type: engine.RB_Kinematic,
		Mass: 1,
		Vx:   0,
		Vy:   0,
	})

	// Add Ball
	ball := engine.AddEntity(world)
	engine.AddComponent(world, ball, &ballComponent{})
	engine.AddComponent(
		world,
		ball,
		&engine.TransformComponent{
			X: float64(engine.Window.CanvasWidth / 2),
			Y: float64(engine.Window.CanvasHeight / 2),
		},
	)
	engine.AddComponent(world, ball, &engine.BoxCollider2DComponent{
		T: 5, B: 5, L: 5, R: 5,
		Category:     1,
		CategoryMask: 1,
	})
	engine.AddComponent(world, ball, &engine.RigidBody2DComponent{
		Type: engine.RB_Dynamic,
		Mass: 1,
		Vx:   -100,
		Vy:   -100,
	})

	// Add Walls
	topWall := engine.AddEntity(world)
	engine.AddComponent(world, topWall, &engine.TransformComponent{
		X: float64(engine.Window.CanvasWidth / 2),
		Y: 0,
	})
	engine.AddComponent(world, topWall, &engine.BoxCollider2DComponent{
		T: 5, B: 5, L: float64(engine.Window.CanvasWidth / 2), R: float64(engine.Window.CanvasWidth / 2),
		Category:     1,
		CategoryMask: 1,
	})
	engine.AddComponent(world, topWall, &engine.RigidBody2DComponent{
		Type: engine.RB_Static,
		Mass: 1,
		Vx:   0,
		Vy:   0,
	})

	bottomWall := engine.AddEntity(world)
	engine.AddComponent(world, bottomWall, &engine.TransformComponent{
		X: float64(engine.Window.CanvasWidth / 2),
		Y: float64(engine.Window.CanvasHeight),
	})
	engine.AddComponent(world, bottomWall, &engine.BoxCollider2DComponent{
		T: 5, B: 5, L: float64(engine.Window.CanvasWidth / 2), R: float64(engine.Window.CanvasWidth / 2),
		Category:     1,
		CategoryMask: 1,
	})
	engine.AddComponent(world, bottomWall, &engine.RigidBody2DComponent{
		Type: engine.RB_Static,
		Mass: 1,
		Vx:   0,
		Vy:   0,
	})
}

// OnExit implements engine.IScene.
func (p *pongScene) OnExit(world *engine.World) {
	engine.Input.Unregister(p.playerInputContext)
	engine.Systems.Unregister(p.systems...)
}

var _ engine.IScene = (*pongScene)(nil)

func NewPongScene() *pongScene {
	playerController := newPlayerController()
	return &pongScene{
		systems: []engine.System{
			playerController,
			// newGameplaySystem(),
			engine.NewPhysics2DSystem(),
			newPongRenderingSystem(),
		},
		playerInputContext: newPlayerInputContext(),
		playerController:   playerController,
	}
}

// ------------------ Player Input Context ------------------

const (
	action_MoveUp engine.InputAction = iota
	action_MoveDown
)

func newPlayerInputContext() *engine.InputMappingContext {
	return engine.NewInputMappingContext().
		RegisterAction(action_MoveUp, engine.NewInputTrigger(engine.Triggered, engine.KeyW, false)).
		RegisterAction(action_MoveDown, engine.NewInputTrigger(engine.Triggered, engine.KeyS, false))
}

// ------------------ Player Controller ------------------

type playerController struct {
	entity engine.Entity
}

// Update implements engine.System.
func (p *playerController) Update(world *engine.World, delta float64) {
	if p.entity.IsNull() {
		return
	}
	transform := engine.GetComponent[engine.TransformComponent](world, p.entity)
	if transform == nil {
		return
	}

	if engine.Input.Get(action_MoveUp) > 0 {
		transform.Y -= 1
	} else if engine.Input.Get(action_MoveDown) > 0 {
		transform.Y += 1
	}
}

func (p *playerController) SetEntity(entity engine.Entity) {
	p.entity = entity
}

var _ engine.System = (*playerController)(nil)

func newPlayerController() *playerController {
	return &playerController{
		entity: engine.Entity{},
	}
}

// ------------------ Components ------------------

type (
	paddleComponent struct{}
	ballComponent   struct{}
)

// ------------------ Gameplay System ------------------
type gameplaySystem struct{}

// Update implements engine.System.
func (g *gameplaySystem) Update(world *engine.World, dt float64) {
	// TODO: Implement ball movement? Physics?
	// TODO: Implement score + reset
}

var _ engine.System = (*gameplaySystem)(nil)

func newGameplaySystem() *gameplaySystem {
	return &gameplaySystem{}
}

// ------------------ Rendering System ------------------
type pongRenderingSystem struct{}

// Update implements engine.System.
func (g *pongRenderingSystem) Update(world *engine.World, dt float64) {
	engine.Query3(
		world,
		func(_ engine.Entity, trans *engine.TransformComponent, collider *engine.BoxCollider2DComponent, _ *paddleComponent) {
			renderer.AddCall(0, 0, func() {
				// TODO: Abstract RayLib calls behind engine/render api package
				rl.DrawRectangle(
					int32(trans.X-collider.L),
					int32(trans.Y-collider.T),
					int32(collider.L+collider.R),
					int32(collider.T+collider.B),
					rl.White,
				)
			})
		},
	)
	engine.Query3(
		world,
		func(_ engine.Entity, trans *engine.TransformComponent, collider *engine.BoxCollider2DComponent, _ *ballComponent) {
			renderer.AddCall(0, 0, func() {
				rl.DrawRectangle(
					int32(trans.X-collider.L),
					int32(trans.Y-collider.T),
					int32(collider.L+collider.R),
					int32(collider.T+collider.B),
					rl.White,
				)
			})
		},
	)
}

var _ engine.System = (*pongRenderingSystem)(nil)

func newPongRenderingSystem() *pongRenderingSystem {
	return &pongRenderingSystem{}
}
