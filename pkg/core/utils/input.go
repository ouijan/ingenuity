package utils

type InputAction HashId

type InputActionMap[T any] = map[InputAction]T

type InputActionValues = InputActionMap[float64]

func GetAction(iav InputActionValues, action InputAction) float64 {
	value, ok := iav[action]
	if !ok {
		return 0
	}
	return value
}

var INPUT_UP = InputAction(Hash("input.up"))
var INPUT_DOWN = InputAction(Hash("input.down"))
var INPUT_LEFT = InputAction(Hash("input.left"))
var INPUT_RIGHT = InputAction(Hash("input.right"))
var INPUT_ZOOM_IN = InputAction(Hash("input.zoomIn"))
var INPUT_ZOOM_OUT = InputAction(Hash("input.zoomOut"))
var INPUT_BOOST = InputAction(Hash("input.boost"))
