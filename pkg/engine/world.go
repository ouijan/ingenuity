package engine

import (
	"fmt"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/ouijan/aether/pkg/core"
)

type IWorld struct {
	ecs ecs.World
}

func NewWorld() *IWorld {
	ecsWorld := ecs.NewWorld()
	ecsWorld.SetListener(newWorldEventProxy())
	return &IWorld{
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

func AddEntity(world *IWorld) Entity {
	return newEntity(world.ecs.NewEntity())
}

func RemoveEntity(world *IWorld, entity Entity) {
	world.ecs.RemoveEntity(entity.entity)
}

func AddComponent[T Component](world *IWorld, entity Entity, component *T) {
	mapper := generic.NewMap1[T](&world.ecs)
	mapper.Assign(entity.entity, component)
}

func RemoveComponent[T Component](world *IWorld, entity Entity, component *T) {
	mapper := generic.NewMap1[T](&world.ecs)
	mapper.Remove(entity.entity)
}

func AddParent(world *IWorld, child Entity, parent Entity) {
	AddComponent(world, child, &ChildOf{Parent: parent.entity})
}

func RemoveParent(world *IWorld, child Entity) {
	RemoveComponent(world, child, &ChildOf{})
}

var World = NewWorld()

func DebugWorld(world *IWorld) {
	core.Log.Debug(fmt.Sprintf("World Stats: \n%s", world.ecs.Stats().String()))
}
