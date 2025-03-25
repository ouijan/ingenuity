package engine

import (
	"encoding/json"
	"fmt"
	"reflect"

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

func (w *World) preSolve(arb *cp.Arbiter) bool {
	a, b := arb.Shapes()
	aData, aOk := a.UserData.(physicsData)
	bData, bOk := b.UserData.(physicsData)
	if !aOk || !bOk {
		core.Log.Error("Failed to get physics data from shape")
		return true
	}

	aEvent := CollisionEvent{TargetEntity: aData.e, OtherEntity: bData.e}
	aData.c.Collisions = append(aData.c.Collisions, aEvent)
	core.EmitEvent("engine.physics.collision", aEvent)

	bEvent := CollisionEvent{TargetEntity: bData.e, OtherEntity: aData.e}
	bData.c.Collisions = append(bData.c.Collisions, bEvent)
	core.EmitEvent("engine.physics.collision", bEvent)

	return aData.rb != nil && bData.rb != nil
}

func (w *World) PrintDebug() {
	msg := "World Debug:\n"
	msg += fmt.Sprintf("\n%s\n", w.ecs.Stats().String())

	comps := make([]ecs.CompInfo, 0)
	for _, id := range ecs.ComponentIDs(&w.ecs) {
		info, ok := ecs.ComponentInfo(&w.ecs, id)
		if ok {
			comps = append(comps, info)
		}
	}

	dump := w.ecs.DumpEntities()
	for _, entity := range dump.Entities {
		if entity.IsZero() {
			continue
		}

		msg += fmt.Sprintf("Entity %d:\n", entity.ID())

		for _, comp := range comps {
			// mapper := generic.NewMap[comp.Type](&w.ecs)
			pointer := w.ecs.Get(entity, comp.ID)

			if pointer != nil {
				value := reflect.NewAt(comp.Type, pointer)
				b, _ := json.Marshal(value.Interface())
				msg += fmt.Sprintf(
					"  %s: %s \n",
					comp.Type.Name(),
					b,
				)
			}
		}

	}

	// Add Physics debug?

	core.Log.Debug(msg)
}

func NewWorld() *World {
	w := &World{
		ecs:    ecs.NewWorld(),
		space:  cp.NewSpace(),
		bodies: map[uint32]*cp.Body{},
		shapes: map[uint32]*cp.Shape{},
	}

	w.ecs.SetListener(newWorldEventProxy(w))
	w.space.Iterations = 5

	handler := w.space.NewWildcardCollisionHandler(cp.WILDCARD_COLLISION_TYPE)
	handler.PreSolveFunc = func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		return w.preSolve(arb)
	}

	return w
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

type Resource interface{}

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

func AddEntity(w *World) Entity {
	return newEntity(w.ecs.NewEntity())
}

func RemoveEntity(w *World, e Entity) {
	w.ecs.RemoveEntity(e.entity)
}

func AddComponent[T Component](w *World, e Entity, comp *T) {
	mapper := generic.NewMap1[T](&w.ecs)
	mapper.Assign(e.entity, comp)
}

func RemoveComponent[T Component](w *World, e Entity, comp *T) {
	mapper := generic.NewMap1[T](&w.ecs)
	mapper.Remove(e.entity)
}

func GetComponent[T Component](w *World, entity Entity) *T {
	mapper := generic.NewMap1[T](&w.ecs)
	return mapper.Get(entity.entity)
}

func AddResource[T Resource](w *World, resource *T) {
	ecs.AddResource(&w.ecs, resource)
}

func GetResource[T Resource](w *World) *T {
	gridRes := generic.NewResource[T](&w.ecs)
	return gridRes.Get()
}

func AddParent(w *World, child Entity, parent Entity) {
	AddComponent(w, child, &ChildOf{Parent: parent.entity})
}

func RemoveParent(w *World, child Entity) {
	RemoveComponent(w, child, &ChildOf{})
}

func Query1[C1 any](w *World, iterator func(Entity, *C1)) {
	filter := generic.NewFilter1[C1]()
	query := filter.Query(&w.ecs)
	for query.Next() {
		c1 := query.Get()
		iterator(newEntity(query.Entity()), c1)
	}
}

func Query2[C1 any, C2 any](w *World, iterator func(Entity, *C1, *C2)) {
	filter := generic.NewFilter2[C1, C2]()
	query := filter.Query(&w.ecs)
	for query.Next() {
		c1, c2 := query.Get()
		iterator(newEntity(query.Entity()), c1, c2)
	}
}

func Query3[C1 any, C2 any, C3 any](w *World, iterator func(Entity, *C1, *C2, *C3)) {
	filter := generic.NewFilter3[C1, C2, C3]()
	query := filter.Query(&w.ecs)
	for query.Next() {
		c1, c2, c3 := query.Get()
		iterator(newEntity(query.Entity()), c1, c2, c3)
	}
}

func Query4[C1 any, C2 any, C3 any, C4 any](
	w *World,
	iterator func(Entity, *C1, *C2, *C3, *C4),
) {
	filter := generic.NewFilter4[C1, C2, C3, C4]()
	query := filter.Query(&w.ecs)
	for query.Next() {
		c1, c2, c3, c4 := query.Get()
		iterator(newEntity(query.Entity()), c1, c2, c3, c4)
	}
}

var CurrentWorld = NewWorld()
