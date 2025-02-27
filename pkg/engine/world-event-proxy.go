package engine

import (
	"fmt"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
	"github.com/ouijan/aether/pkg/core"
)

type worldEventProxy struct{}

func (l *worldEventProxy) Subscriptions() event.Subscription {
	return event.All
}

func (l *worldEventProxy) Components() *ecs.Mask {
	return nil
}

func (l *worldEventProxy) Notify(world *ecs.World, evt ecs.EntityEvent) {
	emitEvent(evt, event.EntityCreated, "entityCreated")
	emitEvent(evt, event.EntityRemoved, "entityRemoved")
	emitEvent(evt, event.ComponentAdded, "componentAdded")
	emitEvent(evt, event.ComponentRemoved, "componentRemoved")
	emitEvent(evt, event.RelationChanged, "relationChanged")
	emitEvent(evt, event.TargetChanged, "targetChanged")
}

func emitEvent(evt ecs.EntityEvent, eventType event.Subscription, eventName string) {
	if evt.EventTypes.Contains(eventType) {
		eventId := fmt.Sprintf("engine.world.%s", eventName)
		core.EmitEvent(eventId, evt)
	}
}

func newWorldEventProxy() *worldEventProxy {
	return &worldEventProxy{}
}
