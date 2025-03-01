package renderer

import (
	"image"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ouijan/ingenuity/pkg/core"
	"github.com/ouijan/ingenuity/pkg/resources"
)

type TilemapLayerRenderCall struct {
	Layer resources.TilemapLayer
}

func NewTilemapLayerRenderCall(
	tilemap *resources.Tilemap,
	layer resources.TilemapLayer,
) TilemapLayerRenderCall {
	return TilemapLayerRenderCall{
		Layer: layer,
	}
}

func RenderTilemapLayer(call TilemapLayerRenderCall) error {
	for tileIndex, tile := range call.Layer.Tiles {
		err := renderTile(call.Layer, tile, tileIndex)
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
) error {
	if tile.IsNil {
		return nil
	}

	texture, ok := TextureCache.Get(cacheId(tile.TilesetID, strconv.Itoa(tile.ID)))
	if !ok {
		return nil
	}

	x, y := layer.GetTilePosition(tileIndex)
	rl.DrawTexture(texture, int32(x), int32(y), rl.White)
	return nil
}

func cacheId(segments ...string) string {
	return strings.Join(segments, ":")
}

func LoadTilemapTextures(tilemap resources.Tilemap) {
	for _, tileset := range tilemap.Tilesets {
		LoadTilesetTextures(tileset)
	}
}

func LoadTilesetTextures(tileset resources.Tileset) {
	tilesetImg := rl.LoadImage(tileset.ImgSrc)
	if tilesetImg == nil {
		core.Log.Error("Failed to load tileset image: " + tileset.ImgSrc)
		return
	}

	for id := 0; id < tileset.TileCount; id++ {
		tileImg := rl.ImageCopy(tilesetImg)

		rl.ImageCrop(tileImg, toRaylibRect(tileset.GetTileRect(id)))
		texture := rl.LoadTextureFromImage(tileImg)
		rl.UnloadImage(tileImg)

		TextureCache.Set(cacheId(tileset.ID, strconv.Itoa(id)), texture)
	}

	rl.UnloadImage(tilesetImg)
}

func toRaylibRect(rect image.Rectangle) rl.Rectangle {
	return rl.NewRectangle(
		float32(rect.Min.X),
		float32(rect.Min.Y),
		float32(rect.Dx()),
		float32(rect.Dy()),
	)
}

func UnloadTilemapTextures(tilemap resources.Tilemap) {
	for _, tileset := range tilemap.Tilesets {
		for id := 0; id < tileset.TileCount; id++ {
			cacheKey := cacheId(tileset.ID, strconv.Itoa(id))
			if texture, ok := TextureCache.Get(cacheKey); ok {
				rl.UnloadTexture(texture)
				TextureCache.Clear(cacheKey)
			}
		}
	}
}
