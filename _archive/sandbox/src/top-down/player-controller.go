package top_down

import "github.com/ouijan/ingenuity/pkg/engine"

type PlayerController struct {
	entity engine.Entity
}

// Update implements engine.System.
func (p *PlayerController) Update(w *engine.World, _ float64) {
	if p.entity.IsNull() {
		return
	}
	transform := engine.GetComponent[engine.TransformComponent](w, p.entity)
	if transform == nil {
		return
	}

	if engine.Input.Get(Player_MoveUp) > 0 {
		transform.Y -= 1
	} else if engine.Input.Get(Player_MoveDown) > 0 {
		transform.Y += 1
	}
	if engine.Input.Get(Player_MoveLeft) > 0 {
		transform.X -= 1
	} else if engine.Input.Get(Player_MoveRight) > 0 {
		transform.X += 1
	}
}

func (p *PlayerController) SetEntity(entity engine.Entity) {
	p.entity = entity
}

var _ engine.System = (*PlayerController)(nil)

func NewPlayerController() PlayerController {
	return PlayerController{}
}
