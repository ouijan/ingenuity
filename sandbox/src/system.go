package main

import (
	"github.com/ouijan/ingenuity/pkg/engine"
)

type SandboxSystem struct{}

var _ engine.System = (*SandboxSystem)(nil)

func (system *SandboxSystem) Update(_ *engine.World) {
	// core.EmitEvent("sandbox.sandboxSystem.update", nil)
	// world.GetEntityCount()
	// Log.Info("SandboxSystem Update")
}

func NewSandboxSystem() *SandboxSystem {
	return &SandboxSystem{}
}
