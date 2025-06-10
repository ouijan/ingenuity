package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ark "github.com/mlange-42/ark/ecs"

	"github.com/ouijan/ingenuity/pkg/client/resources"
	"github.com/ouijan/ingenuity/pkg/core/ecs"
	"github.com/ouijan/ingenuity/pkg/core/utils"
)

type InputHandlerSystem struct {
	store  ark.Resource[resources.UserInputStore]
	camera ark.Resource[rl.Camera2D]
}

// OnCreate implements ecs.System.
func (s *InputHandlerSystem) OnCreate(ea *ecs.EntityAdmin) {
	s.store = ark.NewResource[resources.UserInputStore](&ea.World)
	s.camera = ark.NewResource[rl.Camera2D](&ea.World)
}

// OnDestroy implements ecs.System.
func (i *InputHandlerSystem) OnDestroy() {
	// panic("unimplemented")
}

// Update implements ecs.System.
func (s *InputHandlerSystem) Update(dt float32) error {
	store := s.store.Get()
	camera := s.camera.Get()

	moveSpeed := float32(1000)
	zoomSpeed := float32(.1)

	if store == nil || camera == nil {
		return nil
	}

	if utils.GetAction(store.Values, utils.INPUT_BOOST) > 0 {
		moveSpeed *= 2
		zoomSpeed *= 2
	}

	if utils.GetAction(store.Values, utils.INPUT_ZOOM_IN) > 0 {
		camera.Zoom += zoomSpeed * dt
	} else if utils.GetAction(store.Values, utils.INPUT_ZOOM_OUT) > 0 {
		camera.Zoom -= zoomSpeed * dt
	}

	moveSpeed /= camera.Zoom

	if utils.GetAction(store.Values, utils.INPUT_UP) > 0 {
		camera.Target.Y -= moveSpeed * dt
	} else if utils.GetAction(store.Values, utils.INPUT_DOWN) > 0 {
		camera.Target.Y += moveSpeed * dt
	}
	if utils.GetAction(store.Values, utils.INPUT_LEFT) > 0 {
		camera.Target.X -= moveSpeed * dt
	} else if utils.GetAction(store.Values, utils.INPUT_RIGHT) > 0 {
		camera.Target.X += moveSpeed * dt
	}

	return nil
}

var _ ecs.System = &InputSystem{}
