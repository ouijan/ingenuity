package engine

import (
	"github.com/lafriks/go-tiled"
	"github.com/ouijan/ingenuity/pkg/resources"
)

type TilemapRendererComponent struct {
	TilemapRef *resources.Tilemap
}

type TilemapLayerComponent struct {
	Map   *resources.Tilemap
	Layer resources.TilemapLayer
}

type TiledObjectComponent struct {
	Object *tiled.Object
}

func AddTilemapToWorld(tilemap *resources.Tilemap, world *World) {
	addTilemapLayersToWorld(world, tilemap, tilemap.Layers)
	// addObjectGroupsToWorld(world, tilemap, ObjectGroups)
	// addGroupsToWorld(world, tilemap, tilemap.Tilemap.Groups)
	// renderer.LoadTilemapTextures(tilemap.Tilemap)
}

func addGroupsToWorld(world *World, tilemap *resources.Tilemap, groups []*tiled.Group) {
	for _, group := range groups {
		// addTilemapLayersToWorld(world, tilemap, group.Layers)
		addObjectGroupsToWorld(world, group.ObjectGroups)
		addGroupsToWorld(world, tilemap, group.Groups)
	}
}

func addObjectGroupsToWorld(world *World, groups []*tiled.ObjectGroup) {
	for _, group := range groups {
		for _, obj := range group.Objects {
			addObjectToWorld(world, obj)
		}
	}
}

func addTilemapLayersToWorld(
	world *World,
	tilemap *resources.Tilemap,
	layers []resources.TilemapLayer,
) {
	for _, layer := range layers {
		addLayerToWorld(world, tilemap, layer)
	}
}

func addLayerToWorld(world *World, tilemap *resources.Tilemap, layer resources.TilemapLayer) {
	e := AddEntity(world)
	AddComponent(world, e, &TilemapLayerComponent{
		Layer: layer,
		Map:   tilemap,
	})
}

func addObjectToWorld(world *World, obj *tiled.Object) {
	e := AddEntity(world)
	AddComponent(world, e, &TiledObjectComponent{Object: obj})
}
