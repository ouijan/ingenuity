package input

import "github.com/ouijan/ingenuity/pkg/core/utils"

type InputBindingType uint

const (
	Started InputBindingType = iota
	Completed
	Triggered
)

type inputBinding interface {
	GetValue() float64
}

type inputActionMapping struct {
	Action   utils.InputAction
	Bindings []inputBinding
}

func (ia *inputActionMapping) GetValue() float64 {
	for _, binding := range ia.Bindings {
		value := binding.GetValue()
		if value != 0 {
			return value
		}
	}
	return 0
}
