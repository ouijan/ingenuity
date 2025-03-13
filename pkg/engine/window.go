package engine

import rl "github.com/gen2brain/raylib-go/raylib"

type WindowManager struct {
	CanvasHeight, CanvasWidth int
}

func NewWindowManager() *WindowManager {
	return &WindowManager{}
}

var Window = NewWindowManager()

func (wm *WindowManager) Open(w, h int, title string) {
	// rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagVsyncHint)
	rl.InitWindow(int32(w), int32(h), title)
	rl.SetTargetFPS(60)

	wm.SetCanvasSize(w, h)
}

func (wm *WindowManager) Close() {
	rl.CloseWindow()
}

func (wm *WindowManager) SetCanvasSize(w, h int) {
	wm.CanvasWidth = w
	wm.CanvasHeight = h
	rl.SetWindowMinSize(w, h)
}
