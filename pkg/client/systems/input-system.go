package systems

import (
	"errors"

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

	store.Values = store.Manager.GetAll()
	return nil
}

var _ ecs.System = &InputSystem{}
