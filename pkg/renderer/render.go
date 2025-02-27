package renderer

import (
	"errors"
	"image/png"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Render() error {
	return nil
}

func RenderImage() error {
	tilesetImg := rl.LoadImage("rpg-example/beach_tileset.png")
	if tilesetImg == nil {
		return errors.New("failed to load image")
	}

	tileImg := rl.ImageCopy(tilesetImg)

	// tilesetTexture := rl.LoadTextureFromImage(tilesetImg)
	rl.UnloadImage(tilesetImg)

	var tileX, tileY, tileWidth, tileHeight float32 = 1, 1, 16, 16
	rl.ImageCrop(tileImg, rl.NewRectangle(tileX*tileWidth, tileY*tileHeight, tileWidth, tileHeight))
	tileTexture := rl.LoadTextureFromImage(tileImg)

	xPos, yPos := int32(0), int32(0)
	rl.DrawTexture(
		tileTexture,
		xPos,
		yPos,
		rl.White,
	)
	return nil
}

func openImagePng(path string) (*rl.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	image, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	img := rl.NewImageFromImage(image)
	return img, nil
}
