package engine

import "github.com/ouijan/ingenuity/pkg/core"

var Systems = NewSystemManager()

type System interface {
	Update(world *IWorld)
}

type SystemManager struct {
	systems []System
}

func (sm *SystemManager) Register(system System) {
	sm.systems = append(sm.systems, system)
}

func (sm *SystemManager) Update(world *IWorld) {
	if world == nil {
		core.Log.Error("World not set")
		return
	}
	for _, system := range sm.systems {
		system.Update(world)
	}
}

func NewSystemManager() *SystemManager {
	return &SystemManager{}
}
