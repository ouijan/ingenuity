package renderer

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ouijan/ingenuity/pkg/resources"
)

func RenderTilemapLayer(layer resources.TilemapLayer, xOffset, yOffset int) error {
	for tileIndex, tile := range layer.Tiles {
		err := renderTile(layer, tile, tileIndex, xOffset, yOffset)
		if err != nil {
			return err
		}
	}
	return nil
}

func renderTile(
	layer resources.TilemapLayer,
	tile resources.TilemapTile,
	tileIndex int,
	xOffset int,
	yOffset int,
) error {
	if tile.IsNil {
		return nil
	}

	texture, ok := TextureCache.Get(cacheId(tile.TilesetID, strconv.Itoa(tile.ID)))
	if !ok {
		return nil
	}

	xTile, yTile := layer.GetTilePosition(tileIndex)
	x := xTile + xOffset
	y := yTile + yOffset
	rl.DrawTexture(texture, int32(x), int32(y), rl.White)
	return nil
}

func LoadTilemapTextures(tilemap resources.Tilemap) {
	for _, tileset := range tilemap.Tilesets {
		LoadSpriteSheetTextures(*tileset.GetSpriteSheet())
	}
}

func UnloadTilemapTextures(tilemap resources.Tilemap) {
	for _, tileset := range tilemap.Tilesets {
		UnloadSpritesheetTextures(*tileset.GetSpriteSheet())
	}
}
