package systems

import (
	"fmt"
	// _ "image/png" // Register PNG decoder
	"path"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
	ark "github.com/mlange-42/ark/ecs"

	"github.com/ouijan/ingenuity/pkg/client/renderer"
	"github.com/ouijan/ingenuity/pkg/core/ecs"
	"github.com/ouijan/ingenuity/pkg/core/ecs/components"
	"github.com/ouijan/ingenuity/pkg/core/log"
	"github.com/ouijan/ingenuity/pkg/core/resources"
	"github.com/ouijan/ingenuity/pkg/core/utils"
)

type TinkerSystem struct {
	textEntities  *ark.Filter3[components.Metadata, components.Transform2D, components.Text]
	netEntities   *ark.Filter1[components.NetworkedEntity]
	netComponents *ark.Map3[components.Metadata, components.Transform2D, components.Text]
	tilemap       *tiled.Map
	rm            *resources.ResourceManager
}

func (s *TinkerSystem) OnCreate(ea *ecs.EntityAdmin) {
	s.rm = resources.NewResourceManager("./")

	s.textEntities = ark.NewFilter3[components.Metadata, components.Transform2D, components.Text](&ea.World)
	s.netEntities = ark.NewFilter1[components.NetworkedEntity](&ea.World)
	s.netComponents = ark.NewMap3[components.Metadata, components.Transform2D, components.Text](
		&ea.World,
	)

	tilemapPath := "assets/Tiled Sample/Sample.tmx"
	tilemap, err := s.rm.TiledMap.Load(tilemapPath)
	if err != nil {
		panic(err)
	}

	LoadTilemapTextures(s.rm, tilemap, tilemapPath)
	s.tilemap = tilemap
}

func (s *TinkerSystem) Update(dt float32) error {
	// Tilemap Render

	for i, layer := range s.tilemap.Layers {
		DrawLayer(s.tilemap, layer, i)
	}

	// Text Update + REnder
	fpsMsg := fmt.Sprintf("FPS: %d, DT: %f", rl.GetFPS(), rl.GetFrameTime())

	query := s.textEntities.Query()
	for query.Next() {
		meta, trans, text := query.Get()

		switch meta.Name {
		case utils.FPSDisplayName:
			text.Content = fpsMsg
		}

		renderer.Call(0, 0.0, func() {
			rl.DrawText(
				text.Content,
				int32(trans.X),
				int32(trans.Y),
				int32(text.FontSize),
				rl.LightGray,
			)
		})
	}

	// imgPath := "assets/Tiled Sample/Sprites/crate_E.png"
	// crateImg := rl.LoadImage(path.Join(s.rm.Root, imgPath))
	// crateTexture := rl.LoadTextureFromImage(crateImg)
	// rl.UnloadImage(crateImg)

	// renderer.Call(0, 0.0, func() {

	// 	rl.DrawTexture(
	// 		crateTexture,
	// 		0, 0, rl.White,
	// 	)
	// })

	return nil
}

func (s *TinkerSystem) OnDestroy() {
	// Cleanup if needed
}

var _ ecs.System = &TinkerSystem{}

// Testing

func LoadTilemapTextures(rm *resources.ResourceManager, tilemap *tiled.Map, tilemapPath string) {
	tilemapDir := path.Dir(tilemapPath)

	for _, tileset := range tilemap.Tilesets {
		var tilesetImg *rl.Image
		if tileset.Image != nil {
			tilesetImagePath := path.Join(rm.Root, tilemapDir, tileset.Image.Source)
			tilesetImg = rl.LoadImage(tilesetImagePath)
		}

		for _, tile := range tileset.Tiles {
			// Use sprite image, if available, otherwise use tileset image
			var spriteImg *rl.Image
			if tile.Image != nil {
				spriteImagePath := path.Join(tilemapDir, tile.Image.Source)
				spriteImg = rl.LoadImage(spriteImagePath)
			} else if tilesetImg != nil {
				spriteImg = rl.ImageCopy(tilesetImg)
				rl.ImageCrop(spriteImg, rl.NewRectangle(
					float32(tile.X),
					float32(tile.Y),
					float32(tile.Width),
					float32(tile.Height),
				))
			}

			if spriteImg == nil {
				log.Error("Error loading tile image: %v", tile.Image.Source)
				continue
			}
			texture := rl.LoadTextureFromImage(spriteImg)
			rl.UnloadImage(spriteImg)

			textureId := renderer.TextureID(tileset.Source, fmt.Sprint(tile.ID))
			log.Debug("Loading texture for tile %d in tileset %s: ", tile.ID, tileset.Source, textureId)
			renderer.TextureCache.Set(textureId, texture)
		}

		if tilesetImg != nil {
			rl.UnloadImage(tilesetImg)
		}
	}

}

func UnloadTilemapTextures(tilemap *tiled.Map) {
	for _, tileset := range tilemap.Tilesets {
		for _, tile := range tileset.Tiles {
			textureId := renderer.TextureID(tileset.Source, fmt.Sprint(tile.ID))
			if texture, ok := renderer.TextureCache.Get(textureId); ok {
				rl.UnloadTexture(texture)
				renderer.TextureCache.Clear(textureId)
			}
		}
	}
}

func DrawLayer(tilemap *tiled.Map, layer *tiled.Layer, layerIndex int) {
	if !layer.Visible {
		return
	}

	for tileIndex, tile := range layer.Tiles {
		if tile.IsNil() {
			continue
		}
		textureId := renderer.TextureID(tile.Tileset.Source, fmt.Sprint(tile.ID))
		texture, ok := renderer.TextureCache.Get(textureId)
		if !ok {
			fmt.Printf("Texture not found for tile %d in layer %s\n", tile.ID, layer.Name)
			continue
		}

		xTile := tileIndex % tilemap.Width
		yTile := tileIndex / tilemap.Width
		isoX, isoY := utils.ToIsoCoordinates(xTile, yTile, tilemap.TileWidth, tilemap.TileHeight)
		originX, originY := toTileOriginCoordinates(isoX, isoY, tilemap.TileWidth, tilemap.TileHeight)
		x, y := applySpriteOffset(originX, originY, tilemap.TileHeight, int(texture.Height))

		x += layer.OffsetX
		y += layer.OffsetY

		// labelPadding := 2
		// labelSize := 10
		// labelX := originX + labelPadding
		// labelY := originY + tilemap.TileHeight - labelSize - labelPadding

		renderer.Call(layerIndex, float64(tileIndex), func() {
			rl.DrawTexture(texture, int32(x), int32(y), rl.White)
			// Texture Debug
			// rl.DrawRectangleLines(
			// 	int32(x),
			// 	int32(y),
			// 	int32(texture.Width),
			// 	int32(texture.Height),
			// 	rl.Fade(rl.Blue, 0.5),
			// )
			// Debug Origin
			// rl.DrawRectangleLines(
			// 	int32(originX),
			// 	int32(originY),
			// 	int32(tilemap.TileWidth),
			// 	int32(tilemap.TileHeight),
			// 	rl.Fade(rl.Blue, 0.5),
			// )
			// rl.DrawText(
			// 	fmt.Sprintf("%d, %d (%dpx, %dpx)", xTile, yTile, x, y),
			// 	int32(labelX),
			// 	int32(labelY),
			// 	int32(labelSize),
			// 	rl.Red,
			// )
		})

	}
}

// applySpriteOffset adjusts the position to account for the size of the sprite's
// texture. For tiles with a texture that is larger than the tile size. Uses the
// bottom left corner of the texture as the tiles bottom left corner
func applySpriteOffset(x, y, tileH, textureH int) (int, int) {
	spriteYOffset := textureH - tileH
	originY := y - spriteYOffset
	return x, originY
}

// toTileOriginCoordinates adjusts the origin of the sprite to the center of the tile
func toTileOriginCoordinates(x, y, tileW, tileH int) (int, int) {
	originX := x - tileW/2
	originY := y - tileH/2
	return originX, originY
}
