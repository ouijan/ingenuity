package input

import (
	"slices"

	"github.com/ouijan/ingenuity/pkg/core/utils"
)

type InputManager struct {
	contexts []*InputContext
}

func (im *InputManager) Register(contexts ...*InputContext) {
	im.contexts = append(im.contexts, contexts...)
}

func (im *InputManager) Unregister(contexts ...*InputContext) {
	for _, context := range contexts {
		index := slices.Index(im.contexts, context)
		if index >= 0 {
			im.contexts = slices.Delete(im.contexts, index, index)
		}
	}
}

func (im *InputManager) Get(action utils.InputAction) float64 {
	for _, context := range im.contexts {
		value := context.Get(action)
		if value != 0 {
			return value
		}
	}
	return 0
}

func (im *InputManager) GetAll() utils.InputActionValues {
	values := make(utils.InputActionMap[float64])
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
