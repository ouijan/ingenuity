package pong

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ouijan/ingenuity/pkg/engine"
	"github.com/ouijan/ingenuity/pkg/renderer"
	"github.com/ouijan/ingenuity/sandbox/src/shared"
)

// ------------------ Scene ------------------

type pongScene struct {
	systems            []engine.System
	playerInputContext *engine.InputMappingContext
	playerController   *playerController
	enemyController    *enemyController
	gameplaySystem     *gameplaySystem
}

// Load implements engine.IScene.
func (p *pongScene) Load() {
	// panic("unimplemented")
}

// OnEnter implements engine.IScene.
func (p *pongScene) OnEnter(w *engine.World) {
	// Set Player Controller and Register Systems
	engine.Input.Register(p.playerInputContext)
	engine.Systems.Register(p.systems...)

	halfWidth := float64(engine.Window.CanvasWidth / 2)
	halfHeight := float64(engine.Window.CanvasHeight / 2)

	// Add Player
	player := engine.AddEntity(w)
	engine.AddComponent(w, player, &paddleComponent{})
	engine.AddComponent(w, player, &engine.TransformComponent{
		X: 20, Y: halfHeight,
	})
	engine.AddComponent(w, player, &engine.BoxCollider2DComponent{
		T: 30, B: 30, L: 5, R: 5,
		Category:     1,
		CategoryMask: 1,
	})
	engine.AddComponent(w, player, &engine.RigidBody2DComponent{
		Type: engine.RB_Kinematic,
		Mass: 1,
		Vx:   0,
		Vy:   0,
	})
	p.playerController.SetEntity(player)

	// Add Enemy
	enemy := engine.AddEntity(w)
	engine.AddComponent(w, enemy, &paddleComponent{})
	engine.AddComponent(w, enemy, &engine.TransformComponent{
		X: float64(engine.Window.CanvasWidth - 20),
		Y: halfHeight,
	})
	engine.AddComponent(w, enemy, &engine.BoxCollider2DComponent{
		T: 30, B: 30, L: 5, R: 5,
		Category:     1,
		CategoryMask: 1,
	})
	engine.AddComponent(w, enemy, &engine.RigidBody2DComponent{
		Type: engine.RB_Kinematic,
		Mass: 1,
		Vx:   0,
		Vy:   0,
	})
	p.enemyController.SetEntity(enemy)

	// Add Ball
	ball := engine.AddEntity(w)
	engine.AddComponent(w, ball, &ballComponent{})
	engine.AddComponent(
		w,
		ball,
		&engine.TransformComponent{X: halfWidth, Y: halfHeight},
	)
	engine.AddComponent(w, ball, &engine.BoxCollider2DComponent{
		T: 5, B: 5, L: 5, R: 5,
		Category:     1,
		CategoryMask: 1,
	})
	engine.AddComponent(w, ball, &engine.RigidBody2DComponent{
		Type: engine.RB_Dynamic,
		Mass: 1,
		Vx:   -100,
		Vy:   -100,
	})

	// Add Walls
	topWall := engine.AddEntity(w)
	engine.AddComponent(w, topWall, &engine.TransformComponent{
		X: halfWidth,
		Y: 0,
	})
	engine.AddComponent(w, topWall, &engine.BoxCollider2DComponent{
		T: 5, B: 5, L: halfWidth, R: halfWidth,
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
		X: halfWidth,
		Y: float64(engine.Window.CanvasHeight),
	})
	engine.AddComponent(w, bottomWall, &engine.BoxCollider2DComponent{
		T: 5, B: 5, L: halfWidth, R: halfWidth,
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
	engine.AddComponent(w, leftWall, &engine.TransformComponent{
		X: 0, Y: halfHeight,
	})
	engine.AddComponent(w, leftWall, &engine.BoxCollider2DComponent{
		T: halfHeight, B: halfHeight, L: 5, R: 5,
		Category:     1,
		CategoryMask: 1,
	})
	p.gameplaySystem.LeftGoal = leftWall

	rightWall := engine.AddEntity(w)
	engine.AddComponent(w, rightWall, &engine.TransformComponent{
		X: float64(engine.Window.CanvasWidth), Y: halfHeight,
	})
	engine.AddComponent(w, rightWall, &engine.BoxCollider2DComponent{
		T: halfHeight, B: halfHeight, L: 5, R: 5,
		Category:     1,
		CategoryMask: 1,
	})
	p.gameplaySystem.RightGoal = rightWall

	w.PrintDebug()
}

// OnExit implements engine.IScene.
func (p *pongScene) OnExit(w *engine.World) {
	engine.Input.Unregister(p.playerInputContext)
	engine.Systems.Unregister(p.systems...)
}

var _ engine.IScene = (*pongScene)(nil)

func NewPongScene() *pongScene {
	playerController := newPlayerController()
	enemyController := newEnemyController()
	gameplaySystem := newGameplaySystem()

	return &pongScene{
		systems: []engine.System{
			playerController,
			enemyController,
			engine.NewPhysics2DSystem(),
			gameplaySystem,
			newPongRenderingSystem(),
		},
		playerInputContext: newPlayerInputContext(),
		playerController:   playerController,
		enemyController:    enemyController,
		gameplaySystem:     gameplaySystem,
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

// ------------------ Enemy Controller ------------------

type enemyController struct {
	entity engine.Entity
}

// Update implements engine.System.
func (e *enemyController) Update(world *engine.World, delta float64) {
	if e.entity.IsNull() {
		return
	}
	t := engine.GetComponent[engine.TransformComponent](world, e.entity)
	if t == nil {
		return
	}

	b := e.findBall(world)
	if b == nil {
		return
	}

	if t.Y > b.Y {
		t.Y -= 1
	}
	if t.Y < b.Y {
		t.Y += 1
	}
}

func (e *enemyController) findBall(w *engine.World) *engine.TransformComponent {
	var bt *engine.TransformComponent
	engine.Query2(w, func(_ engine.Entity, t *engine.TransformComponent, _ *ballComponent) {
		bt = t
	})
	return bt
}

func (e *enemyController) SetEntity(entity engine.Entity) {
	e.entity = entity
}

var _ engine.System = (*enemyController)(nil)

func newEnemyController() *enemyController {
	return &enemyController{
		entity: engine.Entity{},
	}
}

// ------------------ Components ------------------

type (
	paddleComponent struct{}
	ballComponent   struct{}
)

// ------------------ Gameplay System ------------------
type gameplaySystem struct {
	LeftGoal   engine.Entity
	RightGoal  engine.Entity
	LeftScore  int
	RightScore int
}

// Update implements engine.System.
func (g *gameplaySystem) Update(w *engine.World, dt float64) {
	shared.Log.Info(fmt.Sprintf("Player 1: %d, Player 2: %d", g.LeftScore, g.RightScore))
	if g.hasBallCollision(w, g.RightGoal) {
		g.LeftScore++
		shared.Log.Info(fmt.Sprintf("Player 1 Scored: %d", g.LeftScore))
		g.resetBall(w)
	}
	if g.hasBallCollision(w, g.LeftGoal) {
		g.RightScore++
		shared.Log.Info(fmt.Sprintf("Player 2 Scored: %d", g.RightScore))
		g.resetBall(w)
	}
}

func (g *gameplaySystem) hasBallCollision(
	w *engine.World,
	e engine.Entity,
) bool {
	if e.IsNull() {
		return false
	}

	col := engine.GetComponent[engine.BoxCollider2DComponent](w, e)
	for _, collision := range col.Collisions {
		ball := engine.GetComponent[ballComponent](w, collision.OtherEntity)
		if ball != nil {
			return true
		}
	}
	return false
}

func (g *gameplaySystem) resetBall(w *engine.World) {
	engine.Query3(
		w,
		func(_ engine.Entity, t *engine.TransformComponent, rb *engine.RigidBody2DComponent, _ *ballComponent) {
			t.X = float64(engine.Window.CanvasWidth / 2)
			t.Y = float64(engine.Window.CanvasHeight / 2)
			rb.Vx = -100
			rb.Vy = -100
		},
	)
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
