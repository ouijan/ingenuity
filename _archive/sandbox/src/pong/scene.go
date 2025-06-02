package pong

import (
	"fmt"
	"path"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ouijan/ingenuity/pkg/audio"
	"github.com/ouijan/ingenuity/pkg/core"
	"github.com/ouijan/ingenuity/pkg/engine"
	"github.com/ouijan/ingenuity/pkg/net"
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

	// TODO: Establish connection with server world (entity admin)
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
	engine.AddComponent(w, player, &PlayerControllerComponent{
		OwnerId: "player1",
	})
	engine.AddComponent(w, player, &PaddleComponent{
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
	engine.AddComponent(w, enemy, &PaddleComponent{
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
	engine.AddComponent(w, ball, &BallComponent{
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
	engine.AddComponent(w, leftWall, &GoalComponent{Id: GoalLeft})
	engine.AddComponent(w, leftWall, &engine.TransformComponent{
		X: 0, Y: halfHeight,
	})
	engine.AddComponent(w, leftWall, &engine.BoxCollider2DComponent{
		T: halfHeight, B: halfHeight, L: 5, R: 5,
		Category:     1,
		CategoryMask: 1,
	})

	rightWall := engine.AddEntity(w)
	engine.AddComponent(w, rightWall, &GoalComponent{Id: GoalRight})
	engine.AddComponent(w, rightWall, &engine.TransformComponent{
		X: float64(engine.Window.CanvasWidth), Y: halfHeight,
	})
	engine.AddComponent(w, rightWall, &engine.BoxCollider2DComponent{
		T: halfHeight, B: halfHeight, L: 5, R: 5,
		Category:     1,
		CategoryMask: 1,
	})

	engine.AddResource(w, &GameState{
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
			NewPlayerController(),
			newEnemyController(),
			engine.NewPhysics2DSystem(),
			NewGameplaySystem(),
			newPongRenderingSystem(),
		},
		playerInputContext: NewPlayerInputContext(),
	}
}

// ------------------ Player Input Context ------------------

const (
	Action_MoveUp engine.InputAction = iota
	Action_MoveDown
)

func NewPlayerInputContext() *engine.InputMappingContext {
	return engine.NewInputMappingContext().
		RegisterAction(Action_MoveUp, engine.NewInputTrigger(engine.Triggered, engine.KeyW, false)).
		RegisterAction(Action_MoveDown, engine.NewInputTrigger(engine.Triggered, engine.KeyS, false))
}

// ------------------ Player Controller ------------------

type playerController struct{}

// Update implements engine.System.
func (pc *playerController) Update(w *engine.World, delta float64) {
	engine.Query3(
		w,
		func(_ engine.Entity, paddle *PaddleComponent, t *engine.TransformComponent, pcComp *PlayerControllerComponent) {
			// Resolve player input for this player
			actions, ok := pc.GetInputFromPlayerId(pcComp.OwnerId)
			if !ok {
				return
			}

			if engine.GetAction(actions, Action_MoveUp) > 0 {
				t.Y -= paddle.Speed
			} else if engine.GetAction(actions, Action_MoveDown) > 0 {
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

func NewPlayerController() *playerController {
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

	engine.Query2(w, func(e engine.Entity, t *engine.TransformComponent, p *PaddleComponent) {
		pc := engine.GetComponent[PlayerControllerComponent](w, e)
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
	engine.Query2(w, func(_ engine.Entity, t *engine.TransformComponent, _ *BallComponent) {
		bt = t
	})
	return bt
}

var _ engine.System = (*enemyController)(nil)

func newEnemyController() *enemyController {
	return &enemyController{}
}

// ------------------ Components ------------------
type PlayerControllerComponent struct {
	OwnerId string
}

type PaddleComponent struct {
	Speed float64
}

type BallComponent struct {
	Speed float64
}

type GoalId uint8

const (
	GoalLeft GoalId = iota
	GoalRight
)

type GoalComponent struct {
	Id GoalId
}

// ------------------ Resources ------------------

type GameState struct {
	LeftScore   int
	RightScore  int
	bounceSound audio.Sound
	winSound    audio.Sound
	loseSound   audio.Sound
}

// ------------------ Gameplay System ------------------
type GameplaySystem struct{}

// Update implements engine.System.
func (g *GameplaySystem) Update(w *engine.World, dt float64) {
	gs := engine.GetResource[GameState](w)
	if gs == nil {
		return
	}

	winSound := soundMap["win"]
	loseSound := soundMap["lose"]
	bounceSound := soundMap["bounce"]

	engine.Query2(w, func(_ engine.Entity, b *BallComponent, col *engine.BoxCollider2DComponent) {
		for _, collision := range col.Collisions {
			goal := engine.GetComponent[GoalComponent](w, collision.OtherEntity)
			if goal != nil {
				if goal.Id == GoalRight {
					engine.Audio.PlaySound(winSound)
					gs.LeftScore++
					shared.Log.Info(fmt.Sprintf("Player 1 Scored: %d", gs.LeftScore))
					g.resetBall(w)
					return
				}
				if goal.Id == GoalLeft {
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

func (g *GameplaySystem) resetBall(w *engine.World) {
	engine.Query3(
		w,
		func(_ engine.Entity, t *engine.TransformComponent, rb *engine.RigidBody2DComponent, b *BallComponent) {
			t.X = float64(engine.Window.CanvasWidth / 2)
			t.Y = float64(engine.Window.CanvasHeight / 2)
			rb.Vx = -100 * b.Speed
			rb.Vy = -100 * b.Speed
		},
	)
}

var _ engine.System = (*GameplaySystem)(nil)

func NewGameplaySystem() *GameplaySystem {
	return &GameplaySystem{}
}

// ------------------ Rendering System ------------------
type pongRenderingSystem struct{}

// Update implements engine.System.
func (g *pongRenderingSystem) Update(w *engine.World, dt float64) {
	engine.Query3(
		w,
		func(_ engine.Entity, trans *engine.TransformComponent, collider *engine.BoxCollider2DComponent, _ *PaddleComponent) {
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
		func(_ engine.Entity, trans *engine.TransformComponent, collider *engine.BoxCollider2DComponent, _ *BallComponent) {
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

	gs := engine.GetResource[GameState](w)
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

// ------------------ Player Spawn System ------------------

type playerSpawnSystem struct {
	connectedCh    chan net.Addr
	disconnectedCh chan net.Addr
}

// Update implements engine.System.
func (p *playerSpawnSystem) Update(w *engine.World, dt float64) {
	err := core.ReadCh(p.connectedCh, func(addr net.Addr) error {
		return p.handlePlayerConnected(w, addr)
	})
	if err != nil {
		core.Log.Error(err.Error())
	}
	err = core.ReadCh(p.disconnectedCh, func(addr net.Addr) error {
		return p.handlePlayerDisconnected(w, addr)
	})
	if err != nil {
		core.Log.Error(err.Error())
	}
}

func (p *playerSpawnSystem) handlePlayerConnected(
	w *engine.World,
	addr net.Addr,
) error {
	var existing *PlayerControllerComponent
	var vacant *PlayerControllerComponent

	engine.Query1(w, func(e engine.Entity, pc *PlayerControllerComponent) {
		if pc.OwnerId == addr.String() {
			existing = pc
		}
		if pc.OwnerId == "" {
			vacant = pc
		}
	})

	if existing != nil {
		core.Log.Info(fmt.Sprintf("Player already connected: %s", addr.String()))
		return nil
	}
	if vacant != nil {
		vacant.OwnerId = addr.String()
		core.Log.Info(fmt.Sprintf("Player connected: %s", addr.String()))
		return nil
	}
	core.Log.Info(fmt.Sprintf("Player limit reached: %s", addr.String()))
	// core.EmitEvent("engine.network.disconnectUser", addr)
	return nil
}

func (p *playerSpawnSystem) handlePlayerDisconnected(w *engine.World, addr net.Addr) error {
	engine.Query1(w, func(e engine.Entity, pc *PlayerControllerComponent) {
		if pc.OwnerId == addr.String() {
			pc.OwnerId = ""
		}
	})
	core.Log.Info(fmt.Sprintf("Player diconnected: %s", addr.String()))
	return nil
}

var _ engine.System = (*playerSpawnSystem)(nil)

func NewPlayerSpawnSystem() *playerSpawnSystem {
	p := &playerSpawnSystem{
		connectedCh:    make(chan net.Addr, 2),
		disconnectedCh: make(chan net.Addr, 2),
	}
	core.OnEventCh("engine.network.userConnected", p.connectedCh)
	core.OnEventCh("engine.network.userDisconnected", p.disconnectedCh)
	return p
}
