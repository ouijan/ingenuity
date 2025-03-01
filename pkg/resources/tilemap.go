package resources

import (
	"fmt"
	"image"

	"github.com/lafriks/go-tiled"
	"github.com/ouijan/ingenuity/pkg/core"
)

// Tilemap

type Tilemap struct {
	Layers   []TilemapLayer
	Tilesets []Tileset
	// We should be able to construct a tilemap without using tiled-go
	// We should provide a tiled-integration module that builds our tilemap from a go-tiled map
}

// Tilemap Layer

type TilemapLayer struct {
	Name       string
	TileWidth  int
	TileHeight int
	Width      int
	Tiles      []TilemapTile
}

func (tl *TilemapLayer) GetTilePosition(tileID int) (int, int) {
	x := (tileID % tl.Width) * tl.TileWidth
	y := (tileID / tl.Width) * tl.TileHeight
	return x, y
}

func ToTilemapLayer(tilemap *tiled.Map, layer *tiled.Layer) TilemapLayer {
	tiles := []TilemapTile{}

	for _, tile := range layer.Tiles {
		tilesetID := ""
		if tile.Tileset != nil {
			tilesetID = tile.Tileset.Source
		}
		tiles = append(tiles, TilemapTile{
			ID:        int(tile.ID),
			IsNil:     tile.IsNil(),
			TilesetID: tilesetID,
		})
	}
	return TilemapLayer{
		Name:       layer.Name,
		Tiles:      tiles,
		TileWidth:  tilemap.TileWidth,
		TileHeight: tilemap.TileHeight,
		Width:      tilemap.Width,
	}
}

type TilemapTile struct {
	// TileSet *Tileset
	ID        int
	IsNil     bool
	TilesetID string
}

// Tileset

type Tileset struct {
	ID         string
	Name       string
	ImgSrc     string
	TileCount  int
	Columns    int
	TileWidth  int
	TileHeight int
}

func (ts *Tileset) GetTileRect(id int) image.Rectangle {
	x := id % ts.Columns
	y := id / ts.Columns
	return image.Rect(
		x*ts.TileWidth,
		y*ts.TileHeight,
		(x+1)*ts.TileWidth,
		(y+1)*ts.TileHeight,
	)
}

func ToTileset(tileset *tiled.Tileset) Tileset {
	return Tileset{
		ID:         tileset.Source,
		Name:       tileset.Name,
		ImgSrc:     tileset.GetFileFullPath(tileset.Image.Source),
		Columns:    tileset.Columns,
		TileCount:  tileset.TileCount,
		TileWidth:  tileset.TileWidth,
		TileHeight: tileset.TileHeight,
	}
}

// Load

func LoadTilemap(path string) *Tilemap {
	tilemap, err := tiled.LoadFile(path)
	if err != nil {
		core.Log.Error(fmt.Sprintf("Error loading tilemap: %s\n", err.Error()))
	}

	for _, tileset := range tilemap.Tilesets {
		core.Log.Debug(fmt.Sprintf("Tileset: %s", tileset.Name))
	}

	layers := make([]TilemapLayer, 0)
	for _, layer := range tilemap.Layers {
		layers = append(layers, ToTilemapLayer(tilemap, layer))
	}

	tilesets := make([]Tileset, 0)
	for _, tileset := range tilemap.Tilesets {
		tilesets = append(tilesets, ToTileset(tileset))
	}

	return &Tilemap{
		Layers:   layers,
		Tilesets: tilesets,
	}
}
