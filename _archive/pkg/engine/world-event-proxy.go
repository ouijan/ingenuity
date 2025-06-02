package engine

import (
	"fmt"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
	"github.com/ouijan/ingenuity/pkg/core"
)

// World Event
type WorldEvent struct {
	World *World
	Evt   ecs.EntityEvent
}

// World Event Proxy

type worldEventProxy struct {
	world *World
}

func (l *worldEventProxy) Subscriptions() event.Subscription {
	return event.All
}

func (l *worldEventProxy) Components() *ecs.Mask {
	return nil
}

func (l *worldEventProxy) Notify(world *ecs.World, evt ecs.EntityEvent) {
	l.emitEvent(evt, event.EntityCreated, "entityCreated")
	l.emitEvent(evt, event.EntityRemoved, "entityRemoved")
	l.emitEvent(evt, event.ComponentAdded, "componentAdded")
	l.emitEvent(evt, event.ComponentRemoved, "componentRemoved")
	l.emitEvent(evt, event.RelationChanged, "relationChanged")
	l.emitEvent(evt, event.TargetChanged, "targetChanged")
}

func (l *worldEventProxy) emitEvent(
	evt ecs.EntityEvent,
	eventType event.Subscription,
	eventName string,
) {
	if evt.EventTypes.Contains(eventType) {
		eventId := fmt.Sprintf("engine.world.%s", eventName)
		core.EmitEvent(eventId, WorldEvent{
			World: l.world,
			Evt:   evt,
		})
	}
}

func newWorldEventProxy(world *World) *worldEventProxy {
	return &worldEventProxy{world: world}
}
