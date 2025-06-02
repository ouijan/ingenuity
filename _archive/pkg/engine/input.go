package engine

import (
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ouijan/ingenuity/pkg/core"
)

// Input Keys
const (
	KeyW     = rl.KeyW
	KeyA     = rl.KeyA
	KeyS     = rl.KeyS
	KeyD     = rl.KeyD
	KeyQ     = rl.KeyQ
	KeyE     = rl.KeyE
	KeySpace = rl.KeySpace
)

// Input Trigger

type InputTriggerType uint

const (
	Started InputTriggerType = iota
	Completed
	Triggered
)

type InputTrigger struct {
	Type        InputTriggerType
	Key         int32
	InvertValue bool
}

func (ib *InputTrigger) GetValue() float64 {
	rawValue := ib.rawValue()
	if ib.InvertValue {
		return -rawValue
	}
	return rawValue
}

func (ib *InputTrigger) rawValue() float64 {
	switch ib.Type {
	case Started:
		if rl.IsKeyPressed(ib.Key) {
			return 1
		}
		return 0
	case Completed:
		if rl.IsKeyReleased(ib.Key) {
			return 1
		}
		return 0
	case Triggered:
		if rl.IsKeyDown(ib.Key) {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func NewInputTrigger(
	bindingType InputTriggerType,
	key int32,
	invertValue bool,
) *InputTrigger {
	return &InputTrigger{
		Type:        bindingType,
		Key:         key,
		InvertValue: invertValue,
	}
}

// Input Action

type InputAction uint

type inputActionMapping struct {
	Action   InputAction
	Triggers []*InputTrigger
}

func (ia *inputActionMapping) GetValue() float64 {
	for _, trigger := range ia.Triggers {
		value := trigger.GetValue()
		if value != 0 {
			return value
		}
	}
	return 0
}

func newInputActionMapping(action InputAction, triggers ...*InputTrigger) *inputActionMapping {
	return &inputActionMapping{
		Action:   action,
		Triggers: triggers,
	}
}

// Input Context

type InputMappingContext struct {
	actions map[InputAction]*inputActionMapping
}

func NewInputMappingContext() *InputMappingContext {
	return &InputMappingContext{
		actions: map[InputAction]*inputActionMapping{},
	}
}

func (im *InputMappingContext) RegisterAction(
	actionId InputAction,
	triggers ...*InputTrigger,
) *InputMappingContext {
	im.actions[actionId] = newInputActionMapping(
		actionId,
		triggers...,
	)
	return im
}

func (im *InputMappingContext) Get(actionId InputAction) float64 {
	action, ok := im.actions[actionId]
	if !ok {
		return 0
	}
	return action.GetValue()
}

// Input Manager

type inputManager struct {
	contexts []*InputMappingContext
}

func (im *inputManager) Register(contexts ...*InputMappingContext) {
	im.contexts = append(im.contexts, contexts...)
}

func (im *inputManager) Unregister(contexts ...*InputMappingContext) {
	for _, context := range contexts {
		index := slices.Index(im.contexts, context)
		if index >= 0 {
			im.contexts = slices.Delete(im.contexts, index, index)
		}
	}
}

func (im *inputManager) Get(action InputAction) float64 {
	for _, context := range im.contexts {
		value := context.Get(action)
		if value != 0 {
			return value
		}
	}
	return 0
}

func (im *inputManager) GetAll() InputActionValues {
	values := map[InputAction]float64{}
	for _, context := range im.contexts {
		for actionId := range context.actions {
			value := context.Get(actionId)
			if value != 0 {
				values[actionId] = value
			}
		}
	}
	return values
}

func (im *inputManager) Update() {
	core.EmitEvent("client.input", InputActionsEvent{
		Actions: im.GetAll(),
	})
}

func newInputManager() *inputManager {
	return &inputManager{}
}

var Input = newInputManager()

type InputActionValues = map[InputAction]float64

func GetAction(iav InputActionValues, action InputAction) float64 {
	value, ok := iav[action]
	if !ok {
		return 0
	}
	return value
}

type InputActionsEvent struct {
	Actions InputActionValues
}
