package renderer

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ouijan/ingenuity/pkg/core"
	"github.com/ouijan/ingenuity/pkg/resources"
)

func RenderSprite(spritesheet *resources.SpriteSheet, spriteId int, x, y int) error {
	texture, ok := TextureCache.Get(cacheId(spritesheet.ID, strconv.Itoa(spriteId)))
	if !ok {
		return nil
	}

	rl.DrawTexture(texture, int32(x), int32(y), rl.White)
	return nil
}

func LoadSpriteSheetTextures(spritesheet resources.SpriteSheet) {
	sheetImg := rl.LoadImage(spritesheet.Src)
	if sheetImg == nil {
		core.Log.Error("Failed to load spritesheet image: " + spritesheet.Src)
		return
	}

	for _, sprite := range spritesheet.Sprites {
		spriteImg := rl.ImageCopy(sheetImg)

		rl.ImageCrop(spriteImg, core.RL_Rect(sprite.Rect))
		texture := rl.LoadTextureFromImage(spriteImg)
		rl.UnloadImage(spriteImg)

		TextureCache.Set(cacheId(spritesheet.ID, strconv.Itoa(sprite.ID)), texture)
	}

	rl.UnloadImage(sheetImg)
}

func UnloadSpritesheetTextures(spritesheet resources.SpriteSheet) {
	for _, sprite := range spritesheet.Sprites {
		cacheKey := cacheId(spritesheet.ID, strconv.Itoa(sprite.ID))
		if texture, ok := TextureCache.Get(cacheKey); ok {
			rl.UnloadTexture(texture)
			TextureCache.Clear(cacheKey)
		}
	}
}
