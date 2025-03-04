package engine

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/ouijan/ingenuity/pkg/core"
	"github.com/ouijan/ingenuity/pkg/renderer"
)

func Run() {
	defer close()
	attachRaylibLogger(true)
	attachEventLogger()

	// TODO: Pull these from config (or pass config in)
	Window.Open(800, 450, "Ingenuity")

	for !rl.WindowShouldClose() {
		Scene.Update()
		// Controllers.Update(CurrentWorld)
		Systems.Update(CurrentWorld)
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
	renderCalls := make([]renderer.RenderCall, 0)

	// Establish list to render
	tilemapLayerFilter := generic.NewFilter1[TilemapLayerComponent]()
	tilemapLayerQuery := tilemapLayerFilter.Query(&CurrentWorld.ecs)
	for tilemapLayerQuery.Next() {
		tilemapLayer := tilemapLayerQuery.Get()
		renderCalls = append(
			renderCalls,
			func() error {
				return renderer.RenderTilemapLayer(tilemapLayer.Layer, 0, 0)
			},
		)
	}

	spriteFilter := generic.NewFilter2[SpriteRendererComponent, TransformComponent]()
	spriteQuery := spriteFilter.Query(&CurrentWorld.ecs)
	for spriteQuery.Next() {
		sprite, transform := spriteQuery.Get()
		renderCalls = append(
			renderCalls,
			func() error {
				x, y := toScreenPosition(transform)
				return renderer.RenderSprite(
					sprite.SpriteSheet,
					sprite.SpriteIndex,
					x,
					y,
				)
			},
		)
	}

	// Cull things that are off screen
	// Sort into render layers
	// Sort items within layers
	renderer.Render(renderCalls, Window.CanvasWidth, Window.CanvasHeight)
}

func toScreenPosition(transform *TransformComponent) (int, int) {
	return int(transform.X), int(transform.Y)
}
