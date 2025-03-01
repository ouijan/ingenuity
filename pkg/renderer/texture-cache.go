package renderer

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ouijan/ingenuity/pkg/core"
)

type textureCache struct {
	cache map[string]rl.Texture2D
}

func (tc *textureCache) Get(id string) (rl.Texture2D, bool) {
	texture, ok := tc.cache[id]
	if !ok {
		core.Log.Warn(fmt.Sprintf("texture not found in cache: %s", id))
	}
	return texture, ok
}

func (tc *textureCache) Set(id string, texture rl.Texture2D) {
	if _, ok := tc.cache[id]; ok {
		core.Log.Warn(fmt.Sprintf("overriding texture in cache: %s", id))
	}
	tc.cache[id] = texture
}

func (tc *textureCache) Clear(id string) {
	delete(tc.cache, id)
}

func (tc *textureCache) ClearAll() {
	tc.cache = make(map[string]rl.Texture2D)
}

func NewTextureCache() textureCache {
	return textureCache{
		cache: make(map[string]rl.Texture2D),
	}
}

var TextureCache = NewTextureCache()
