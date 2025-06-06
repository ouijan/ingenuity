package ecs

import (
	"github.com/mlange-42/ark/ecs"
)

type EntityAdmin struct {
	ecs.World
	isActive bool
	systems  []System
}

func (ea *EntityAdmin) Activate() {
	ea.isActive = true
}

func (ea *EntityAdmin) Deactivate() {
	ea.isActive = false
}

func (ea *EntityAdmin) Update(dt float32) error {
	if !ea.isActive {
		return nil
	}
	for _, system := range ea.systems {
		if err := system.Update(dt); err != nil {
			return err
		}
	}
	return nil
}

func NewEntityAdmin() *EntityAdmin {
	return &EntityAdmin{
		World:    ecs.NewWorld(),
		isActive: false,
		systems:  make([]System, 0),
	}
}

// -- Can have multiple entity admins, but there should be a single "root" entity admin
