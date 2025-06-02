package top_down

import (
	"github.com/ouijan/ingenuity/pkg/engine"
)

const (
	Player_MoveUp engine.InputAction = iota
	Player_MoveDown
	Player_MoveLeft
	Player_MoveRight
)

func NewPlayerInputContext() *engine.InputMappingContext {
	return engine.NewInputMappingContext().
		RegisterAction(Player_MoveUp, engine.NewInputTrigger(engine.Triggered, engine.KeyW, false)).
		RegisterAction(Player_MoveDown, engine.NewInputTrigger(engine.Triggered, engine.KeyS, false)).
		RegisterAction(Player_MoveLeft, engine.NewInputTrigger(engine.Triggered, engine.KeyA, false)).
		RegisterAction(Player_MoveRight, engine.NewInputTrigger(engine.Triggered, engine.KeyD, false))
}
