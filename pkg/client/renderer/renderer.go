package renderer

import rl "github.com/gen2brain/raylib-go/raylib"

var stack = newRenderStack()

func Call(layer int, sortValue float64, render func()) {
	stack.Add(layer, sortValue, render)
}

func Clear() {
	stack.Clear()
}

func Render(camera rl.Camera2D) {
	stack.Render(camera)
}
