package engine

import (
	"slices"

	"github.com/ouijan/ingenuity/pkg/core"
)

var Systems = NewSystemManager()

type System interface {
	Update(world *World)
}

type SystemManager struct {
	systems []System
}

func (sm *SystemManager) Register(system System) {
	sm.systems = append(sm.systems, system)
}

func (sm *SystemManager) Unregister(system System) {
	i := slices.Index(sm.systems, system)
	if i < 0 {
		core.Log.Error("System not found")
		return
	}
	sm.systems = slices.Delete(sm.systems, i, i)
}

func (sm *SystemManager) Update(world *World) {
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
