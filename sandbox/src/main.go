package main

import (
	"github.com/ouijan/ingenuity/pkg/engine"
	"github.com/ouijan/ingenuity/sandbox/src/pong"
)

func main() {
	// engine.Scene.SetNext(NewDemoScene())
	engine.Scene.SetNext(pong.NewPongScene())
	engine.Run()
}
