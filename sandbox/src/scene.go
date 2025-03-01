package main

import (
	"github.com/ouijan/ingenuity/pkg/engine"
	"github.com/ouijan/ingenuity/pkg/renderer"
	"github.com/ouijan/ingenuity/pkg/resources"
)

type SandboxScene struct {
	tilemap *resources.Tilemap
}

var _ engine.IScene = (*SandboxScene)(nil)

func (scene *SandboxScene) Load() {
	Log.Info("Loading scene")
	scene.tilemap = resources.LoadTilemap("rpg-example/island.tmx")
	renderer.LoadTilemapTextures(*scene.tilemap)
}

func (scene *SandboxScene) OnEnter(world *engine.IWorld) {
	Log.Info("Entering scene")
	engine.AddTilemapToWorld(scene.tilemap, world)
}

func (scene *SandboxScene) OnExit(world *engine.IWorld) {
	Log.Info("Exiting scene")
	renderer.TextureCache.ClearAll()
}

func NewDemoScene() *SandboxScene {
	return &SandboxScene{}
}
