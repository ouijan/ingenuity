package engine

import rl "github.com/gen2brain/raylib-go/raylib"

type WindowManager struct{}

func NewWindowManager() *WindowManager {
	return &WindowManager{}
}

var Window = NewWindowManager()

func (wm *WindowManager) Open(width, height int32, title string) {
	rl.InitWindow(width, height, title)
	rl.SetTargetFPS(60)
}

func (wm *WindowManager) Close() {
	rl.CloseWindow()
}
