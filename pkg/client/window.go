package client

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Window struct {
	CanvasHeight, CanvasWidth int
}

func NewWindow() *Window {
	return &Window{}
}

func (wm *Window) Open(w, h int, title string) {
	rl.SetTraceLogLevel(rl.LogNone)
	rl.SetConfigFlags(rl.FlagVsyncHint)

	rl.InitWindow(int32(w), int32(h), title)
	rl.SetTargetFPS(60)

	wm.SetCanvasSize(w, h)
}

func (wm *Window) Close() {
	rl.CloseWindow()
}

func (wm *Window) SetCanvasSize(w, h int) {
	wm.CanvasWidth = w
	wm.CanvasHeight = h
	rl.SetWindowMinSize(w, h)
}
