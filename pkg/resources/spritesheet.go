package resources

import "image"

type SpriteSheet struct {
	ID      string
	Src     string
	Sprites []Sprite
}

type Sprite struct {
	ID   int
	Rect image.Rectangle
}

func NewSpriteSheet(id, src string, sprites []Sprite) *SpriteSheet {
	return &SpriteSheet{
		ID:      id,
		Src:     src,
		Sprites: sprites,
	}
}

func SpritesFromGrid(spriteWidth, spriteHeight, rows, cols int) []Sprite {
	sprites := make([]Sprite, 0, rows*cols)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			sprites = append(sprites, Sprite{
				ID: y*cols + x,
				Rect: image.Rect(
					x*spriteWidth,
					y*spriteHeight,
					(x+1)*spriteWidth,
					(y+1)*spriteHeight,
				),
			})
		}
	}
	return sprites
}
