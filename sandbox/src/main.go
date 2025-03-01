package main

import (
	"github.com/ouijan/ingenuity/pkg/engine"
)

func main() {
	engine.Systems.Register(NewSandboxSystem())

	scene := NewDemoScene()
	engine.Scene.SetNext(scene)

	Log.Info("Running engine")
	engine.Run()
}
