package input

import rl "github.com/gen2brain/raylib-go/raylib"

type keyBinding struct {
	Type        InputBindingType
	Key         int32
	InvertValue bool
}

func KeyBinding(
	key int32,
	bindingType InputBindingType,
	invertValue bool,
) *keyBinding {
	return &keyBinding{
		Type:        bindingType,
		Key:         key,
		InvertValue: invertValue,
	}
}

var _ inputBinding = (*keyBinding)(nil)

func (ib *keyBinding) GetValue() float64 {
	rawValue := ib.rawValue()
	if ib.InvertValue {
		return -rawValue
	}
	return rawValue
}

func (ib *keyBinding) rawValue() float64 {
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
