package renderer

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ouijan/ingenuity/pkg/core/utils"
)

type renderLayerStack = map[int][]func()
type renderStack struct {
	layers map[int]renderLayerStack
}

// TODO: Layer should denote if it is a world layer or a UI layer.

func newRenderStack() *renderStack {
	return &renderStack{
		layers: make(map[int]renderLayerStack),
	}
}

func (rs *renderStack) Add(layer int, sortValue float64, render func()) {
	if rs.layers == nil {
		rs.layers = make(map[int]renderLayerStack)
	}
	if _, ok := rs.layers[layer]; !ok {
		rs.layers[layer] = make(renderLayerStack)
	}
	rs.layers[layer][int(sortValue)] = append(rs.layers[layer][int(sortValue)], render)
}

func (rs *renderStack) Clear() {
	rs.layers = make(map[int]renderLayerStack)
}

func (rs *renderStack) Render(camera rl.Camera2D) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	rl.BeginMode2D(camera)

	for _, layer := range utils.IterateIntKeyedMap(rs.layers) {
		for _, sortValue := range utils.IterateIntKeyedMap(layer) {
			for _, renderFunc := range sortValue {
				renderFunc()
			}
		}
	}
	rl.EndMode2D()

	// UI needs to be drawn after the 2D mode ends
	rl.DrawText(
		fmt.Sprintf("Center %f, %f", camera.Target.X, camera.Target.Y),
		10, 10, 10, rl.Magenta,
	)

	rl.EndDrawing()
}
