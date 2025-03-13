package engine

import (
	"fmt"

	cp "github.com/jakecoffman/cp/v2"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/ouijan/ingenuity/pkg/core"
)

type World struct {
	ecs    ecs.World
	space  *cp.Space
	bodies map[uint32]*cp.Body
	shapes map[uint32]*cp.Shape
}

func (w *World) Reset() {
	w.ecs.Reset()
	w.bodies = map[uint32]*cp.Body{}
	w.shapes = map[uint32]*cp.Shape{}
	w.space.EachBody(func(body *cp.Body) {
		w.space.RemoveBody(body)
	})
	w.space.EachShape(func(shape *cp.Shape) {
		w.space.RemoveShape(shape)
	})
	w.space.EachConstraint(func(constraint *cp.Constraint) {
		w.space.RemoveConstraint(constraint)
	})
}

func (w *World) cacheBody(e Entity, body *cp.Body) {
	w.bodies[e.ID()] = body
}

func (w *World) getBody(e Entity) *cp.Body {
	body, ok := w.bodies[e.ID()]
	if !ok {
		return nil
	}
	return body
}

func (w *World) cacheShape(e Entity, shape *cp.Shape) {
	w.shapes[e.ID()] = shape
}

func (w *World) getShape(e Entity) *cp.Shape {
	shape, ok := w.shapes[e.ID()]
	if !ok {
		return nil
	}
	return shape
}

func NewWorld() *World {
	world := &World{
		ecs:    ecs.NewWorld(),
		space:  cp.NewSpace(),
		bodies: map[uint32]*cp.Body{},
		shapes: map[uint32]*cp.Shape{},
	}

	world.ecs.SetListener(newWorldEventProxy(world))
	world.space.Iterations = 5

	// Testing - Start

	// xMid := 200.0
	// yMid := 200.0
	// aEntity := AddEntity(world)
	// aTrans := &TransformComponent{X: xMid, Y: yMid + 50}
	// upsertPhysicsEntity(world, aEntity, aTrans, nil, nil)
	// aCol := &BoxCollider2DComponent{T: 15, B: 15, L: 15, R: 15, Category: 1, CategoryMask: 1}
	// upsertPhysicsEntity(world, aEntity, aTrans, aCol, nil)
	// aRb := &RigidBody2DComponent{Type: RB_Dynamic, Mass: 1, Vx: 0, Vy: -50}
	// upsertPhysicsEntity(world, aEntity, aTrans, aCol, aRb)
	//
	// bEntity := AddEntity(world)
	// bTrans := &TransformComponent{X: xMid, Y: yMid - 100}
	// upsertPhysicsEntity(world, bEntity, bTrans, nil, nil)
	// bCol := &BoxCollider2DComponent{T: 15, B: 15, L: 15, R: 15, Category: 1, CategoryMask: 1}
	// upsertPhysicsEntity(world, bEntity, bTrans, bCol, nil)
	// bRb := &RigidBody2DComponent{Type: RB_Dynamic, Mass: 1, Vx: 0, Vy: 50}
	// upsertPhysicsEntity(world, bEntity, bTrans, bCol, bRb)
	//
	DebugWorld(world)

	// Testing - End
	return world
}

// Entity

type Entity struct {
	entity ecs.Entity
}

func (e *Entity) IsNull() bool {
	return e.entity.IsZero()
}

func (e *Entity) ID() uint32 {
	return e.entity.ID()
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

func GetComponent[T Component](world *World, entity Entity) *T {
	mapper := generic.NewMap1[T](&world.ecs)
	return mapper.Get(entity.entity)
}

func AddParent(world *World, child Entity, parent Entity) {
	AddComponent(world, child, &ChildOf{Parent: parent.entity})
}

func RemoveParent(world *World, child Entity) {
	RemoveComponent(world, child, &ChildOf{})
}

func Query1[C1 any](world *World, iterator func(Entity, *C1)) {
	filter := generic.NewFilter1[C1]()
	query := filter.Query(&world.ecs)
	for query.Next() {
		c1 := query.Get()
		iterator(newEntity(query.Entity()), c1)
	}
}

func Query2[C1 any, C2 any](world *World, iterator func(Entity, *C1, *C2)) {
	filter := generic.NewFilter2[C1, C2]()
	query := filter.Query(&world.ecs)
	for query.Next() {
		c1, c2 := query.Get()
		iterator(newEntity(query.Entity()), c1, c2)
	}
}

func Query3[C1 any, C2 any, C3 any](world *World, iterator func(Entity, *C1, *C2, *C3)) {
	filter := generic.NewFilter3[C1, C2, C3]()
	query := filter.Query(&world.ecs)
	for query.Next() {
		c1, c2, c3 := query.Get()
		iterator(newEntity(query.Entity()), c1, c2, c3)
	}
}

func Query4[C1 any, C2 any, C3 any, C4 any](
	world *World,
	iterator func(Entity, *C1, *C2, *C3, *C4),
) {
	filter := generic.NewFilter4[C1, C2, C3, C4]()
	query := filter.Query(&world.ecs)
	for query.Next() {
		c1, c2, c3, c4 := query.Get()
		iterator(newEntity(query.Entity()), c1, c2, c3, c4)
	}
}

var CurrentWorld = NewWorld()

func DebugWorld(world *World) {
	core.Log.Debug(fmt.Sprintf("World Stats: \n%s", world.ecs.Stats().String()))
}
