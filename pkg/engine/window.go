package engine

import rl "github.com/gen2brain/raylib-go/raylib"

type WindowManager struct {
	CanvasHeight, CanvasWidth int
}

func NewWindowManager() *WindowManager {
	return &WindowManager{}
}

var Window = NewWindowManager()

func (wm *WindowManager) Open(width, height int, title string) {
	wm.SetCanvasSize(width, height)
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(int32(width), int32(height), title)
	rl.SetTargetFPS(60)
}

func (wm *WindowManager) Close() {
	rl.CloseWindow()
}

func (wm *WindowManager) SetCanvasSize(width, height int) {
	wm.CanvasWidth = width
	wm.CanvasHeight = height
}
