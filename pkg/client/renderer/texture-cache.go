package renderer

import (
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ouijan/ingenuity/pkg/core/log"
	"github.com/ouijan/ingenuity/pkg/core/utils"
)

type textureCache struct {
	cache map[utils.HashId]rl.Texture2D
}

func (tc *textureCache) Get(id utils.HashId) (rl.Texture2D, bool) {
	texture, ok := tc.cache[id]
	if !ok {
		log.Warn("Texture not found in cache: %v", id)
	}
	return texture, ok
}

func (tc *textureCache) Set(id utils.HashId, texture rl.Texture2D) {
	if _, ok := tc.cache[id]; ok {
		log.Warn("Overriding texture in cache: %v", id)
	}
	tc.cache[id] = texture
}

func (tc *textureCache) Clear(id utils.HashId) {
	delete(tc.cache, id)
}

func (tc *textureCache) ClearAll() {
	tc.cache = make(map[utils.HashId]rl.Texture2D)
}

func NewTextureCache() textureCache {
	return textureCache{
		cache: make(map[utils.HashId]rl.Texture2D),
	}
}

var TextureCache = NewTextureCache()

func TextureID(segments ...string) utils.HashId {
	return utils.Hash(strings.Join(segments, ":"))
}
