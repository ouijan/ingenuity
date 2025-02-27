package engine

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/ouijan/aether/pkg/core"
	"github.com/ouijan/aether/pkg/renderer"
)

func Run() {
	defer close()
	attachRaylibLogger(true)
	attachEventLogger()

	Window.Open(800, 450, "raylib [core] example - basic window")

	for !rl.WindowShouldClose() {
		Scene.Update()
		Systems.Update(World)
		render()
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

func attachEventLogger() {
	core.OnEvent("engine.world.*", func(evt core.Event[ecs.EntityEvent]) error {
		entityId := evt.Data.Entity.ID()
		core.Log.Info(fmt.Sprintf("%s: %v", evt.EventId, entityId))
		return nil
	})
}

func render() {
	// Establish list to render
	tilemapLayerFilter := generic.NewFilter1[TilemapLayerComponent]()
	tilemapLayerQuery := tilemapLayerFilter.Query(&World.ecs)
	renderCalls := make([]renderer.TilemapLayerRenderCall, 0)
	for tilemapLayerQuery.Next() {
		tilemapLayer := tilemapLayerQuery.Get()
		renderCalls = append(
			renderCalls,
			renderer.NewTilemapLayerRenderCall(tilemapLayer.Map, tilemapLayer.Layer),
		)
	}

	// Cull things that are off screen
	// Sort into render layers
	// Sort items in layers

	// Start Drawing
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	// Render layers
	rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

	for _, renderCall := range renderCalls {
		renderer.RenderTilemapLayer(renderCall)
	}

	// err := renderer.RenderImage()
	// if err != nil {
	// 	core.Log.Error(fmt.Sprintf("Error rendering image: %s", err.Error()))
	// }

	// End Drawing
	rl.EndDrawing()
}
