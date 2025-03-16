package engine

import (
	"fmt"
	"os"
	"runtime/trace"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ouijan/ingenuity/pkg/audio"
	"github.com/ouijan/ingenuity/pkg/core"
	"github.com/ouijan/ingenuity/pkg/renderer"
)

func Run() {
	defer close()

	attachRaylibLogger(true)
	attachEventLogger(false, false)
	attachTraceLogger(false)

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

func attachRaylibLogger(quiet bool) {
	if quiet {
		rl.SetTraceLogLevel(rl.LogNone)
		return
	}
	rl.SetTraceLogCallback(func(logType int, text string) {
		switch logType {
		case int(rl.LogDebug):
			core.Log.Debug(text)
			break
		case int(rl.LogInfo):
			core.Log.Info(text)
			break
		case int(rl.LogWarning):
			core.Log.Warn(text)
			break
		case int(rl.LogError):
			core.Log.Error(text)
			break
		default:
			core.Log.Info(text)
			break
		}
	})
}

func attachEventLogger(quiet bool, verbose bool) {
	if quiet {
		return
	}
	core.OnEvent("*", func(evt core.Event[any]) error {
		msg := fmt.Sprintf("%s", evt.EventId)
		if verbose {
			msg += fmt.Sprintf(" -> %#v\n", evt.Data)
		}
		core.Log.Info(msg)
		return nil
	})
}

func attachTraceLogger(quiet bool) {
	if quiet {
		return
	}
	f, _ := os.Create("trace.out")
	trace.Start(f)
}
