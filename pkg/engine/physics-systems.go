package engine

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	cp "github.com/jakecoffman/cp/v2"
	"github.com/mlange-42/arche/ecs"
	"github.com/ouijan/ingenuity/pkg/core"
	"github.com/ouijan/ingenuity/pkg/renderer"
)

type physics2DSystem struct {
	DrawDebug bool
}

type physicsData struct {
	e  Entity
	t  *TransformComponent
	c  *BoxCollider2DComponent
	rb *RigidBody2DComponent
}

// Update implements System.
func (p *physics2DSystem) Update(w *World, dt float64) {
	w.space.EachBody(func(body *cp.Body) {
		updateBody(body)
	})
	w.space.EachShape(func(shape *cp.Shape) {
		updateShape(shape)
	})
	w.space.Step(dt)
	w.space.EachBody(func(body *cp.Body) {
		updateData(body)
	})
	if p.DrawDebug {
		drawDebug(w.space)
	}
}

func updateData(body *cp.Body) {
	data, ok := body.UserData.(physicsData)
	if !ok {
		core.Log.Error("Failed to get physics data from body")
		return
	}
	if data.t != nil {
		data.t.X = body.Position().X
		data.t.Y = body.Position().Y
	}
	if data.rb != nil {
		data.rb.Vx = body.Velocity().X
		data.rb.Vy = body.Velocity().Y
	}
}

func preSolve(arb *cp.Arbiter) bool {
	a, b := arb.Shapes()
	aData, ok := a.UserData.(physicsData)
	if !ok {
		core.Log.Error("Failed to get physics data from shape")
		return true
	}
	bData, ok := b.UserData.(physicsData)
	if !ok {
		core.Log.Error("Failed to get physics data from shape")
		return true
	}

	core.Log.Debug(fmt.Sprintf("Collision between %v and %v", aData, bData))
	// entityA.Collider.AddCollision(core.NewCollisionEvent(entityB.Entity))
	// if entityA.Rigidbody != nil {
	// 	entityA.Rigidbody.AddCollision(core.NewCollisionEvent(entityB.Entity))
	// }
	return true
	// return aData.rb != nil && bData.rb != nil
}

func newPhysicsData(
	e Entity,
	t *TransformComponent,
	c *BoxCollider2DComponent,
	rb *RigidBody2DComponent,
) physicsData {
	return physicsData{
		e:  e,
		t:  t,
		c:  c,
		rb: rb,
	}
}

func newBody(
	e Entity,
	t *TransformComponent,
	c *BoxCollider2DComponent,
	rb *RigidBody2DComponent,
) *cp.Body {
	body := cp.NewBody(0, 0)
	if rb != nil {
		switch rb.Type {
		case RB_Static:
			body.SetType(cp.BODY_STATIC)
			break
		case RB_Kinematic:
			body.SetType(cp.BODY_KINEMATIC)
			break
		case RB_Dynamic:
			body.SetType(cp.BODY_DYNAMIC)
			body.SetMass(rb.Mass)
			body.SetMoment(cp.INFINITY)
			break
		}
	} else {
		body.SetType(cp.BODY_STATIC)
	}
	body.UserData = newPhysicsData(e, t, c, rb)
	return body
}

func updateBody(body *cp.Body) {
	data, ok := body.UserData.(physicsData)
	if !ok {
		core.Log.Error("Failed to get physics data from body")
		return
	}
	if (data.t == nil) || (data.c == nil) {
		return
	}
	body.SetPosition(cp.Vector{X: data.t.X, Y: data.t.Y})
	if data.rb != nil {
		body.SetVelocity(data.rb.Vx, data.rb.Vy)
	}
}

func newShape(
	body *cp.Body,
	e Entity,
	t *TransformComponent,
	c *BoxCollider2DComponent,
	rb *RigidBody2DComponent,
) *cp.Shape {
	bb := cp.NewBB(-c.L, -c.B, c.R, c.T)
	shape := cp.NewBox2(body, bb, 0)
	shape.UserData = newPhysicsData(e, t, c, rb)
	return shape
}

func updateShape(shape *cp.Shape) {
	data, ok := shape.UserData.(physicsData)
	if !ok {
		core.Log.Error("Failed to get physics data from shape")
		return
	}
	if (data.t == nil) || (data.c == nil) {
		return
	}

	shape.SetFriction(.1)
	shape.SetElasticity(1)
	shape.SetCollisionType(cp.WILDCARD_COLLISION_TYPE)
	shape.SetFilter(
		cp.NewShapeFilter(cp.NO_GROUP, uint(data.c.Category), uint(data.c.CategoryMask)),
	)

	if data.rb != nil {
		shape.SetSensor(false)
	} else {
		shape.SetSensor(true)
	}
}

func drawDebug(space *cp.Space) {
	space.EachBody(func(body *cp.Body) {
		colour := rl.Red
		if body.GetType() == cp.BODY_STATIC {
			colour = rl.Green
		} else if body.GetType() == cp.BODY_KINEMATIC {
			colour = rl.Blue
		}

		body.EachShape(func(shape *cp.Shape) {
			renderer.AddCall(5, 0, func() {
				x := shape.BB().L
				y := shape.BB().T
				w := shape.BB().R - shape.BB().L
				h := shape.BB().B - shape.BB().T
				rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), colour)
			})
		})
	})
}

var _ System = (*physics2DSystem)(nil)

func NewPhysics2DSystem() *physics2DSystem {
	system := &physics2DSystem{
		DrawDebug: true,
	}

	handler := CurrentWorld.space.NewWildcardCollisionHandler(cp.WILDCARD_COLLISION_TYPE)
	handler.PreSolveFunc = func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		return preSolve(arb)
	}

	core.OnEvent("engine.world.componentAdded", func(evt core.Event[WorldEvent]) error {
		return handleComponentEvent(evt)
	})
	core.OnEvent("engine.world.componentRemoved", func(evt core.Event[WorldEvent]) error {
		return handleComponentEvent(evt)
	})
	core.OnEvent("engine.world.entityRemoved", func(evt core.Event[WorldEvent]) error {
		removePhysicsEntity(evt.Data.World, newEntity(evt.Data.Evt.Entity))
		return nil
	})

	return system
}

func handleComponentEvent(evt core.Event[WorldEvent]) error {
	entity := newEntity(evt.Data.Evt.Entity)
	transId := ecs.ComponentID[TransformComponent](&evt.Data.World.ecs)
	colId := ecs.ComponentID[BoxCollider2DComponent](&evt.Data.World.ecs)
	rbId := ecs.ComponentID[RigidBody2DComponent](&evt.Data.World.ecs)
	mask := ecs.All(transId, colId, rbId)

	if evt.Data.Evt.Removed.ContainsAny(&mask) || evt.Data.Evt.Added.ContainsAny(&mask) {
		trans := GetComponent[TransformComponent](evt.Data.World, entity)
		col := GetComponent[BoxCollider2DComponent](evt.Data.World, entity)
		rb := GetComponent[RigidBody2DComponent](evt.Data.World, entity)
		rebuildPhysicsEntity(evt.Data.World, entity, trans, col, rb)
	}

	return nil
}

func rebuildPhysicsEntity(
	w *World,
	e Entity,
	t *TransformComponent,
	c *BoxCollider2DComponent,
	rb *RigidBody2DComponent,
) {
	if (t == nil) || (c == nil) {
		return
	}

	oldBody := w.getBody(e)
	if oldBody != nil {
		w.space.RemoveBody(oldBody)
	}
	oldShape := w.getShape(e)
	if oldShape != nil {
		w.space.RemoveShape(oldShape)
	}

	body := newBody(e, t, c, rb)
	updateBody(body)
	shape := newShape(body, e, t, c, rb)
	updateShape(shape)

	w.space.AddBody(body)
	w.space.AddShape(shape)

	w.cacheBody(e, body)
	w.cacheShape(e, shape)
}

func removePhysicsEntity(w *World, e Entity) {
	body := w.getBody(e)
	if body != nil {
		w.space.RemoveBody(body)
		delete(w.bodies, e.ID())
	}

	shape := w.getShape(e)
	if shape != nil {
		w.space.RemoveShape(shape)
		delete(w.shapes, e.ID())
	}
}
