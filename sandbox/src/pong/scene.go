package pong

import (
	"fmt"
	"path"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ouijan/ingenuity/pkg/audio"
	"github.com/ouijan/ingenuity/pkg/engine"
	"github.com/ouijan/ingenuity/pkg/renderer"
	"github.com/ouijan/ingenuity/sandbox/src/shared"
)

// ------------------ Scene ------------------

type pongScene struct {
	systems            []engine.System
	playerInputContext *engine.InputMappingContext
}

var soundMap = map[string]audio.Sound{} // TODO: Replace use of global with a resource manager

// Load implements engine.IScene.
func (p *pongScene) Load() {
	rootDir := "./"
	soundMap["bounce"] = audio.LoadSound(path.Join(rootDir, "assets/audio/bounce.wav"))
	soundMap["lose"] = audio.LoadSound(path.Join(rootDir, "assets/audio/lose.wav"))
	soundMap["win"] = audio.LoadSound(path.Join(rootDir, "assets/audio/win.wav"))
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
	engine.AddComponent(w, player, &playerControllerComponent{
		OwnerId: "player1",
	})
	engine.AddComponent(w, player, &paddleComponent{
		Speed: 2,
	})
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

	// Add Enemy
	enemy := engine.AddEntity(w)
	engine.AddComponent(w, enemy, &paddleComponent{
		Speed: 2,
	})
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

	// Add Ball
	ballSpeed := 2.0
	ball := engine.AddEntity(w)
	engine.AddComponent(w, ball, &ballComponent{
		Speed: ballSpeed,
	})
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
		Vx:   -100 * ballSpeed,
		Vy:   -100 * ballSpeed,
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
	engine.AddComponent(w, leftWall, &goalComponent{Id: goalLeft})
	engine.AddComponent(w, leftWall, &engine.TransformComponent{
		X: 0, Y: halfHeight,
	})
	engine.AddComponent(w, leftWall, &engine.BoxCollider2DComponent{
		T: halfHeight, B: halfHeight, L: 5, R: 5,
		Category:     1,
		CategoryMask: 1,
	})

	rightWall := engine.AddEntity(w)
	engine.AddComponent(w, rightWall, &goalComponent{Id: goalRight})
	engine.AddComponent(w, rightWall, &engine.TransformComponent{
		X: float64(engine.Window.CanvasWidth), Y: halfHeight,
	})
	engine.AddComponent(w, rightWall, &engine.BoxCollider2DComponent{
		T: halfHeight, B: halfHeight, L: 5, R: 5,
		Category:     1,
		CategoryMask: 1,
	})

	engine.AddResource(w, &gameState{
		LeftScore:  0,
		RightScore: 0,
	})
	// w.PrintDebug()
}

// OnExit implements engine.IScene.
func (p *pongScene) OnExit(w *engine.World) {
	engine.Input.Unregister(p.playerInputContext)
	engine.Systems.Unregister(p.systems...)
	soundMap = map[string]audio.Sound{} // Don't like the global here
}

var _ engine.IScene = (*pongScene)(nil)

func NewPongScene() *pongScene {
	return &pongScene{
		systems: []engine.System{
			newPlayerController(),
			newEnemyController(),
			engine.NewPhysics2DSystem(),
			newGameplaySystem(),
			newPongRenderingSystem(),
		},
		playerInputContext: newPlayerInputContext(),
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

type playerController struct{}

// Update implements engine.System.
func (pc *playerController) Update(w *engine.World, delta float64) {
	engine.Query3(
		w,
		func(_ engine.Entity, paddle *paddleComponent, t *engine.TransformComponent, pcComp *playerControllerComponent) {
			// Resolve player input for this player
			actions, ok := pc.GetInputFromPlayerId(pcComp.OwnerId)
			if !ok {
				return
			}

			if engine.GetAction(actions, action_MoveUp) > 0 {
				t.Y -= paddle.Speed
			} else if engine.GetAction(actions, action_MoveDown) > 0 {
				t.Y += paddle.Speed
			}
		},
	)
}

func (pc *playerController) GetInputFromPlayerId(playerId string) (engine.InputActionValues, bool) {
	if playerId != "player1" {
		return engine.InputActionValues{}, false
	}
	return engine.Input.GetAll(), true
}

var _ engine.System = (*playerController)(nil)

func newPlayerController() *playerController {
	return &playerController{}
}

// ------------------ Enemy Controller ------------------

type enemyController struct{}

// Update implements engine.System.
func (ec *enemyController) Update(w *engine.World, delta float64) {
	b := ec.findBall(w)
	if b == nil {
		return
	}

	engine.Query2(w, func(e engine.Entity, t *engine.TransformComponent, p *paddleComponent) {
		pc := engine.GetComponent[playerControllerComponent](w, e)
		if pc != nil {
			return
		}
		if t.Y > b.Y {
			t.Y -= p.Speed
		}
		if t.Y < b.Y {
			t.Y += p.Speed
		}
	})
}

func (ec *enemyController) findBall(w *engine.World) *engine.TransformComponent {
	var bt *engine.TransformComponent
	engine.Query2(w, func(_ engine.Entity, t *engine.TransformComponent, _ *ballComponent) {
		bt = t
	})
	return bt
}

var _ engine.System = (*enemyController)(nil)

func newEnemyController() *enemyController {
	return &enemyController{}
}

// ------------------ Components ------------------
type playerControllerComponent struct {
	OwnerId string
}

type paddleComponent struct {
	Speed float64
}

type ballComponent struct {
	Speed float64
}

type goalId uint8

const (
	goalLeft goalId = iota
	goalRight
)

type goalComponent struct {
	Id goalId
}

// ------------------ Resources ------------------

type gameState struct {
	LeftScore   int
	RightScore  int
	bounceSound audio.Sound
	winSound    audio.Sound
	loseSound   audio.Sound
}

// ------------------ Gameplay System ------------------
type gameplaySystem struct{}

// Update implements engine.System.
func (g *gameplaySystem) Update(w *engine.World, dt float64) {
	gs := engine.GetResource[gameState](w)
	if gs == nil {
		return
	}

	winSound := soundMap["win"]
	loseSound := soundMap["lose"]
	bounceSound := soundMap["bounce"]

	engine.Query2(w, func(_ engine.Entity, b *ballComponent, col *engine.BoxCollider2DComponent) {
		for _, collision := range col.Collisions {
			goal := engine.GetComponent[goalComponent](w, collision.OtherEntity)
			if goal != nil {
				if goal.Id == goalRight {
					engine.Audio.PlaySound(winSound)
					gs.LeftScore++
					shared.Log.Info(fmt.Sprintf("Player 1 Scored: %d", gs.LeftScore))
					g.resetBall(w)
					return
				}
				if goal.Id == goalLeft {
					engine.Audio.PlaySound(loseSound)
					gs.RightScore++
					shared.Log.Info(fmt.Sprintf("Player 2 Scored: %d", gs.RightScore))
					g.resetBall(w)
					return
				}
			}
			engine.Audio.PlaySound(bounceSound)
		}
	})
}

func (g *gameplaySystem) resetBall(w *engine.World) {
	engine.Query3(
		w,
		func(_ engine.Entity, t *engine.TransformComponent, rb *engine.RigidBody2DComponent, b *ballComponent) {
			t.X = float64(engine.Window.CanvasWidth / 2)
			t.Y = float64(engine.Window.CanvasHeight / 2)
			rb.Vx = -100 * b.Speed
			rb.Vy = -100 * b.Speed
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
func (g *pongRenderingSystem) Update(w *engine.World, dt float64) {
	engine.Query3(
		w,
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
		w,
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

	gs := engine.GetResource[gameState](w)
	if gs != nil {
		hw := int32(engine.Window.CanvasWidth / 2)
		size := int32(25)
		gap := int32(10)

		lScore := fmt.Sprintf("%d", gs.LeftScore)
		lOffset := rl.MeasureText(lScore, size) + gap

		rScore := fmt.Sprintf("%d", gs.RightScore)
		rOffset := gap

		renderer.AddCall(0, 0, func() {
			rl.DrawText(lScore, hw-lOffset, 25, size, rl.White)
			rl.DrawText(rScore, hw+rOffset, 25, size, rl.White)
		})
	}
}

var _ engine.System = (*pongRenderingSystem)(nil)

func newPongRenderingSystem() *pongRenderingSystem {
	return &pongRenderingSystem{}
}
