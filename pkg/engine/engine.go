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
	attachEventLogger(true)
	attachTraceLogger(false)

	// TODO: Pull these from config (or pass config in)
	Window.Open(800, 450, "Ingenuity")

	for !rl.WindowShouldClose() {
		Scene.Update()
		Systems.Update(CurrentWorld)
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

func attachEventLogger(quiet bool) {
	if quiet {
		return
	}
	core.OnEvent("engine.world.*", func(evt core.Event[WorldEvent]) error {
		entityId := evt.Data.Evt.Entity.ID()
		core.Log.Info(fmt.Sprintf("%s: %v", evt.EventId, entityId))
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
