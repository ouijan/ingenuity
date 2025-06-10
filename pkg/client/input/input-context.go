package input

import "github.com/ouijan/ingenuity/pkg/core/utils"

type InputContext struct {
	actions utils.InputActionMap[*inputActionMapping]
}

func NewInputContext() *InputContext {
	return &InputContext{
		actions: make(utils.InputActionMap[*inputActionMapping]),
	}
}

func (im *InputContext) RegisterAction(
	actionId utils.InputAction,
	bindings ...inputBinding,
) *InputContext {
	im.actions[actionId] = &inputActionMapping{
		Action:   actionId,
		Bindings: bindings,
	}
	return im
}

func (im *InputContext) Get(actionId utils.InputAction) float64 {
	action, ok := im.actions[actionId]
	if !ok {
		return 0
	}
	return action.GetValue()
}
