package engine

import "github.com/ouijan/ingenuity/pkg/renderer"

func toScreenPosition(transform *TransformComponent) (int, int) {
	return int(transform.X), int(transform.Y)
}

// ---------- TilemapLayerRenderSystem ----------

type tilemapLayerRenderSystem struct{}

// Update implements System.
func (t *tilemapLayerRenderSystem) Update(world *World, delta float64) {
	// TODO: Cull entities that are off screen
	Query1(CurrentWorld, func(entity Entity, tilemapLayer *TilemapLayerComponent) {
		renderer.AddCall(0, 0, func() {
			renderer.RenderTilemapLayer(tilemapLayer.Layer, 0, 0)
		})
	})
}

var _ System = (*tilemapLayerRenderSystem)(nil)

func NewTilemapLayerRenderSystem() *tilemapLayerRenderSystem {
	return &tilemapLayerRenderSystem{}
}

// ---------- SpriteRenderSystem ----------

type spriteRenderSystem struct{}

// Update implements System.
func (s *spriteRenderSystem) Update(world *World, delta float64) {
	// TODO: Cull entities that are off screen
	Query2(
		CurrentWorld,
		func(entity Entity, sprite *SpriteRendererComponent, transform *TransformComponent) {
			x, y := toScreenPosition(transform)
			renderer.AddCall(0, 0, func() {
				renderer.RenderSprite(
					sprite.SpriteSheet,
					sprite.SpriteIndex,
					x,
					y,
				)
			})
		},
	)
}

var _ System = (*spriteRenderSystem)(nil)

func NewSpriteRenderSystem() *spriteRenderSystem {
	return &spriteRenderSystem{}
}
