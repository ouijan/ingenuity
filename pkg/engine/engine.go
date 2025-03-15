package engine

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ouijan/ingenuity/pkg/core"
	"github.com/ouijan/ingenuity/pkg/renderer"
)

func Run() {
	defer close()

	attachRaylibLogger(true)
	attachEventLogger(true)

	// TODO: Pull these from config (or pass config in)
	Window.Open(800, 450, "Ingenuity")

	for !rl.WindowShouldClose() {
		Scene.Update()
		Systems.Update(CurrentWorld)
		renderer.Render(Window.CanvasWidth, Window.CanvasHeight)
	}
}

func close() {
	core.Log.Sync()
	core.CloseEvents()
	Window.Close()
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
