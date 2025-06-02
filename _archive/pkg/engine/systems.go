package engine

import (
	"slices"
	"time"

	"github.com/ouijan/ingenuity/pkg/core"
)

var Systems = NewSystemManager()

type System interface {
	Update(world *World, delta float64)
}

type SystemManager struct {
	lastUpdate time.Time
	systems    []System
}

func (sm *SystemManager) Register(systems ...System) {
	for _, system := range systems {
		sm.systems = append(sm.systems, system)
	}
}

func (sm *SystemManager) Unregister(systems ...System) {
	for _, system := range systems {
		i := slices.Index(sm.systems, system)
		if i < 0 {
			core.Log.Error("System not found")
			return
		}
		sm.systems = slices.Delete(sm.systems, i, i)
	}
}

func (sm *SystemManager) Update(world *World, dt float32) {
	// delta := sm.getDelta()
	if world == nil {
		core.Log.Error("World not set")
		return
	}
	for _, system := range sm.systems {
		system.Update(world, float64(dt))
	}
}

func (sm *SystemManager) getDelta() float64 {
	now := time.Now()
	diff := now.Sub(sm.lastUpdate)
	sm.lastUpdate = now
	return diff.Seconds()
}

func NewSystemManager() *SystemManager {
	return &SystemManager{
		lastUpdate: time.Now(),
	}
}
