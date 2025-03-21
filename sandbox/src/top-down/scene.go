package top_down

import (
	"github.com/ouijan/ingenuity/pkg/engine"
	"github.com/ouijan/ingenuity/pkg/renderer"
	"github.com/ouijan/ingenuity/pkg/resources"
	"github.com/ouijan/ingenuity/sandbox/src/shared"
)

type SandboxScene struct {
	systems           []engine.System
	tilemap           *resources.Tilemap
	playerSpriteSheet *resources.SpriteSheet
	inputCtx          *engine.InputMappingContext
	playerController  PlayerController
}

var _ engine.IScene = (*SandboxScene)(nil)

func (scene *SandboxScene) Load() {
	shared.Log.Info("Loading scene")

	scene.tilemap = resources.LoadTilemap("assets/maps/ortho-map.tmx")
	renderer.LoadTilemapTextures(*scene.tilemap)
	// TODO: texture cache should be given a key to track what requested the texture load
	// This key should then be used to release the texture when the scene is unloaded
	// We need a key for this scene, so that if another scene loads the same texture, it doesn't get unloaded
	// Maintain a map of textures to scenes, and when a scene is unloaded, remove the key from the texture
	// If the texture has no keys, unload it

	// TODO: relative file path handling, currently requires running from the root of the project
	scene.playerSpriteSheet = resources.LoadSpriteSheet("assets/sprites/player.png", 16, 16, 1, 1)
	renderer.LoadSpriteSheetTextures(*scene.playerSpriteSheet)
}

func (scene *SandboxScene) OnEnter(world *engine.World) {
	shared.Log.Info("Entering scene")
	// engine.Controllers.Register(scene.playerController)

	engine.AddTilemapToWorld(scene.tilemap, world)

	// Player Entity
	player := engine.AddEntity(world)
	engine.AddComponent(world, player, &engine.TransformComponent{
		X: 50,
		Y: 50,
	})
	engine.AddComponent(world, player, &engine.SpriteRendererComponent{
		SpriteSheet: scene.playerSpriteSheet,
		SpriteIndex: 0,
	})

	// Camera Entity
	camera := engine.AddEntity(world)
	engine.AddComponent(world, camera, &engine.TransformComponent{
		X: 50,
		Y: 50,
	})
	// engine.AddComponent(world, camera, &engine.CameraComponent{})
	// engine.AddParent(world, camera, player)
	// scene.playerController.SetCamera(&world.Camera)

	scene.playerController.SetEntity(player)
	// TODO: Camera Controller

	engine.Input.Register(scene.inputCtx)
	for _, system := range scene.systems {
		engine.Systems.Register(system)
	}
}

func (scene *SandboxScene) OnExit(world *engine.World) {
	shared.Log.Info("Exiting scene")
	renderer.TextureCache.ClearAll()
	engine.Input.Unregister(scene.inputCtx)
	for _, system := range scene.systems {
		engine.Systems.Unregister(system)
	}
}

func NewDemoScene() *SandboxScene {
	playerCtrl := NewPlayerController()
	return &SandboxScene{
		systems: []engine.System{
			&playerCtrl,
			engine.NewTilemapLayerRenderSystem(),
			engine.NewSpriteRenderSystem(),
		},
		inputCtx:         NewPlayerInputContext(),
		playerController: playerCtrl,
	}
}
