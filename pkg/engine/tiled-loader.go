package engine

import (
	"github.com/lafriks/go-tiled"
	"github.com/ouijan/aether/pkg/resources"
)

type TilemapRendererComponent struct {
	TilemapRef *resources.Tilemap
}

type TilemapLayerComponent struct {
	Map   *tiled.Map
	Layer *tiled.Layer
}

type TiledObjectComponent struct {
	Object *tiled.Object
}

func AddTilemapToWorld(tilemap *resources.Tilemap, world *IWorld) {
	addTilemapLayersToWorld(world, tilemap.Tilemap, tilemap.Tilemap.Layers)
	addObjectGroupsToWorld(world, tilemap.Tilemap.ObjectGroups)
	addGroupsToWorld(world, tilemap.Tilemap, tilemap.Tilemap.Groups)
	// renderer.LoadTilemapTextures(tilemap.Tilemap)
}

func addGroupsToWorld(world *IWorld, tilemap *tiled.Map, groups []*tiled.Group) {
	for _, group := range groups {
		addTilemapLayersToWorld(world, tilemap, group.Layers)
		addObjectGroupsToWorld(world, group.ObjectGroups)
		addGroupsToWorld(world, tilemap, group.Groups)
	}
}

func addObjectGroupsToWorld(world *IWorld, groups []*tiled.ObjectGroup) {
	for _, group := range groups {
		for _, obj := range group.Objects {
			addObjectToWorld(world, obj)
		}
	}
}

func addTilemapLayersToWorld(world *IWorld, tilemap *tiled.Map, layers []*tiled.Layer) {
	for _, layer := range layers {
		addLayerToWorld(world, tilemap, layer)
	}
}

func addLayerToWorld(world *IWorld, tilemap *tiled.Map, layer *tiled.Layer) {
	e := AddEntity(world)
	AddComponent(world, e, &TilemapLayerComponent{
		Layer: layer,
		Map:   tilemap,
	})
}

func addObjectToWorld(world *IWorld, obj *tiled.Object) {
	e := AddEntity(world)
	AddComponent(world, e, &TiledObjectComponent{Object: obj})
}
