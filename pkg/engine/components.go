package engine

import "github.com/ouijan/ingenuity/pkg/resources"

type TransformComponent struct {
	X, Y float64
}

type SpriteRendererComponent struct {
	SpriteSheet *resources.SpriteSheet
	SpriteIndex int
}
