package core

import (
	"image"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func RL_Rect(rect image.Rectangle) rl.Rectangle {
	return rl.NewRectangle(
		float32(rect.Min.X),
		float32(rect.Min.Y),
		float32(rect.Dx()),
		float32(rect.Dy()),
	)
}
