package engine

import (
	"fmt"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/ouijan/ingenuity/pkg/core"
)

type World struct {
	ecs ecs.World
}

func NewWorld() *World {
	ecsWorld := ecs.NewWorld()
	ecsWorld.SetListener(newWorldEventProxy())
	return &World{
		ecs: ecsWorld,
	}
}

// Entity

type Entity struct {
	entity ecs.Entity
}

func newEntity(entity ecs.Entity) Entity {
	return Entity{
		entity: entity,
	}
}

// Component

type Component interface{}

type ChildOf struct {
	ecs.Relation
	Parent ecs.Entity
}

// Entity Methods

// func AddEntity1[T Component](world *IWorld, component *T, parent ...Entity) Entity {
// 	// update ecs
// 	builder := generic.NewMap1[T](&world.ecs, generic.T[ChildOf]())
// 	if len(parent) > 0 {
// 		entity := builder.NewWith(component, parent[0].entity)
// 		return newEntity(entity)
// 	}
// 	entity := builder.NewWith(component)
// 	return newEntity(entity)
// }

func AddEntity(world *World) Entity {
	return newEntity(world.ecs.NewEntity())
}

func RemoveEntity(world *World, entity Entity) {
	world.ecs.RemoveEntity(entity.entity)
}

func AddComponent[T Component](world *World, entity Entity, component *T) {
	mapper := generic.NewMap1[T](&world.ecs)
	mapper.Assign(entity.entity, component)
}

func RemoveComponent[T Component](world *World, entity Entity, component *T) {
	mapper := generic.NewMap1[T](&world.ecs)
	mapper.Remove(entity.entity)
}

func AddParent(world *World, child Entity, parent Entity) {
	AddComponent(world, child, &ChildOf{Parent: parent.entity})
}

func RemoveParent(world *World, child Entity) {
	RemoveComponent(world, child, &ChildOf{})
}

var CurrentWorld = NewWorld()

func DebugWorld(world *World) {
	core.Log.Debug(fmt.Sprintf("World Stats: \n%s", world.ecs.Stats().String()))
}
