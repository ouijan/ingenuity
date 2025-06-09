package tiled

// Tileset represents a tileset in a Tiled .tsx file.
type Tileset struct {
	Version      string `xml:"version,attr" json:"version"`
	TiledVersion string `xml:"tiledversion,attr" json:"tiledversion"`
	Name         string `xml:"name,attr" json:"name"`
	TileWidth    int    `xml:"tilewidth,attr" json:"tilewidth"`
	TileHeight   int    `xml:"tileheight,attr" json:"tileheight"`
	TileCount    int    `xml:"tilecount,attr,omitempty" json:"tilecount,omitempty"`
	Columns      int    `xml:"columns,attr,omitempty" json:"columns,omitempty"`
	Grid         *Grid  `xml:"grid,omitempty" json:"grid,omitempty"`
	Tiles        []Tile `xml:"tile" json:"tiles"`
}

// Grid represents the grid element in a tileset.
type Grid struct {
	Orientation string `xml:"orientation,attr" json:"orientation"`
	Width       int    `xml:"width,attr,omitempty" json:"width,omitempty"`
	Height      int    `xml:"height,attr,omitempty" json:"height,omitempty"`
}

// Tile represents an individual tile in a tileset.
type Tile struct {
	ID    int    `xml:"id,attr" json:"id"`
	Type  string `xml:"type,attr,omitempty" json:"type,omitempty"`
	Image *Image `xml:"image" json:"image"`
}

// Image represents the image element in a tileset.
type Image struct {
	Source string `xml:"source,attr" json:"source"`
	Width  int    `xml:"width,attr,omitempty" json:"width,omitempty"`
	Height int    `xml:"height,attr,omitempty" json:"height,omitempty"`
}
