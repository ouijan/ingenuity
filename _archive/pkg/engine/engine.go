package engine

import (
	"runtime/trace"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ouijan/ingenuity/pkg/audio"
	"github.com/ouijan/ingenuity/pkg/core"
	"github.com/ouijan/ingenuity/pkg/renderer"
)

func Run() {
	defer close()

	core.AttachRaylibLogger(true)
	core.AttachEventLogger(false, false)
	core.AttachTraceLogger(false)

	// TODO: Pull these from config (or pass config in)
	Window.Open(800, 450, "Ingenuity")

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		Scene.Update()

		Input.Update()
		Systems.Update(CurrentWorld, dt)

		renderer.AddCall(5, 0, func() {
			rl.DrawFPS(10, 10)
		})
		renderer.Render(Window.CanvasWidth, Window.CanvasHeight)
	}
}

var Audio = audio.NewAudioDevice()

func close() {
	core.Log.Sync()
	core.CloseEvents()
	Window.Close()
	trace.Stop()
}
