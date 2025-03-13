package renderer

import (
	"cmp"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type renderCall struct {
	Layer     int
	SortValue float64
	Render    func()
}

type renderCallManager struct {
	calls []renderCall
}

func (rcm *renderCallManager) Add(layer int, sortValue float64, render func()) {
	rcm.calls = append(rcm.calls, renderCall{
		Layer:     layer,
		SortValue: sortValue,
		Render:    render,
	})
}

func (rcm *renderCallManager) Execute() {
	slices.SortFunc(rcm.calls, func(a, b renderCall) int {
		if a.Layer == b.Layer {
			return cmp.Compare(a.SortValue, b.SortValue)
		}
		return cmp.Compare(a.Layer, b.Layer)
	})

	for _, renderCall := range rcm.calls {
		renderCall.Render()
	}
}

func (rcm *renderCallManager) Clear() {
	rcm.calls = []renderCall{}
}

func newRenderCallManager() *renderCallManager {
	return &renderCallManager{}
}

var rcm = newRenderCallManager()

func AddCall(layer int, sortValue float64, render func()) {
	rcm.Add(layer, sortValue, render)
}

func Render(
	width, height int,
) error {
	canvas := renderCanvasTexture(width, height)
	drawCanvas(canvas, width, height)
	rl.UnloadRenderTexture(canvas)

	return nil
}

func renderCanvasTexture(
	width, height int,
) rl.RenderTexture2D {
	canvas := rl.LoadRenderTexture(int32(width), int32(height))
	rl.SetTextureFilter(canvas.Texture, rl.TextureFilterLinear)

	rl.BeginTextureMode(canvas)
	rl.ClearBackground(rl.DarkGray)
	rcm.Execute()
	rl.EndTextureMode()

	rcm.Clear()

	return canvas
}

func drawCanvas(canvas rl.RenderTexture2D, width, height int) {
	screenWidth := rl.GetScreenWidth()
	screenHeight := rl.GetScreenHeight()
	scale := min(screenWidth/width, screenHeight/height)
	scale = max(1, scale)

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawTexturePro(
		canvas.Texture,
		rl.Rectangle{Width: float32(canvas.Texture.Width), Height: -float32(canvas.Texture.Height)},
		rl.Rectangle{
			X:      float32(screenWidth-(width*scale)) * 0.5,
			Y:      float32(screenHeight-(height*scale)) * 0.5,
			Width:  float32(canvas.Texture.Width) * float32(scale),
			Height: float32(canvas.Texture.Height) * float32(scale),
		},
		rl.Vector2{},
		0,
		rl.White,
	)
	rl.EndDrawing()
}
