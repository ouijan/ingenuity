package systems

import (
	"errors"

	rl "github.com/gen2brain/raylib-go/raylib"
	ark "github.com/mlange-42/ark/ecs"

	"github.com/ouijan/ingenuity/pkg/client/resources"
	"github.com/ouijan/ingenuity/pkg/core/ecs"
)

type InputSystem struct {
	store ark.Resource[resources.UserInputStore]
}

// OnCreate implements ecs.System.
func (s *InputSystem) OnCreate(ea *ecs.EntityAdmin) {
	s.store = ark.NewResource[resources.UserInputStore](&ea.World)
}

// OnDestroy implements ecs.System.
func (i *InputSystem) OnDestroy() {
	// No resources to clean up in this system
	// This method is required by the ecs.System interface
	// but does not need to do anything in this case.
}

// Update implements ecs.System.
func (s *InputSystem) Update(dt float32) error {
	store := s.store.Get()
	if store == nil {
		return errors.New("UserInputStore resource not found")
	}

	store.Up = rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp)
	store.Down = rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown)
	store.Left = rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft)
	store.Right = rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight)
	store.ZoomIn = rl.IsKeyDown(rl.KeyI) || rl.IsKeyDown(rl.KeyEqual)
	store.ZoomOut = rl.IsKeyDown(rl.KeyO) || rl.IsKeyDown(rl.KeyMinus)
	store.Boost = rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)

	return nil
}

var _ ecs.System = &InputSystem{}
