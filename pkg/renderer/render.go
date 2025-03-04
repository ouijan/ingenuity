package renderer

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type RenderCall func() error

func Render(
	renderCalls []RenderCall,
	width, height int,
) error {
	canvas := renderCanvasTexture(renderCalls, width, height)
	drawCanvas(canvas, width, height)
	rl.UnloadRenderTexture(canvas)

	return nil
}

func renderCanvasTexture(
	renderCalls []RenderCall,
	width, height int,
) rl.RenderTexture2D {
	canvas := rl.LoadRenderTexture(int32(width), int32(height))
	rl.SetTextureFilter(canvas.Texture, rl.TextureFilterLinear)

	rl.BeginTextureMode(canvas)
	rl.ClearBackground(rl.White)
	for _, renderCall := range renderCalls {
		renderCall()
	}
	rl.EndTextureMode()

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
