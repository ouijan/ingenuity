package ecs

import (
	"slices"
)

type System interface {
	OnCreate(ea *EntityAdmin)
	Update(dt float32) error
	OnDestroy()
}

func AddSystem(ea *EntityAdmin, systems ...System) {
	for _, system := range systems {
		system.OnCreate(ea)
	}
	ea.systems = append(ea.systems, systems...)
}

func RemoveSystem(ea *EntityAdmin, systems ...System) {
	for _, system := range systems {
		i := slices.Index(ea.systems, system)
		if i >= 0 {
			ea.systems = slices.Delete(ea.systems, i, i)
		}
		system.OnDestroy()
	}
}
