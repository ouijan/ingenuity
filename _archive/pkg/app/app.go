package app

import (
	"slices"
	"time"

	"github.com/ouijan/ingenuity/pkg/core"
)

type Layer interface {
	OnAttach()
	OnDetach()
	OnUpdate(dt float32) error
	// OnRender()
}

type Config struct {
	Name string
}

type App struct {
	config     Config
	layerStack []Layer
	exit       bool
}

func (a *App) Attach(layer Layer) {
	a.layerStack = append(a.layerStack, layer)
	layer.OnAttach()
}

func (a *App) Detach(layer Layer) {
	a.layerStack = slices.DeleteFunc(a.layerStack, func(l Layer) bool {
		if l == layer {
			layer.OnDetach()
			return true
		}
		return false
	})
}

func (a *App) Update(dt float32) error {
	for _, layer := range a.layerStack {
		err := layer.OnUpdate(dt)
		if err != nil {
			return err
		}

	}
	return nil
}

func (a *App) Run(updatesPerSecond int) error {
	core.OnEvent("app.exit", func(e core.Event[any]) error {
		core.Log.Info("Exiting")
		a.exit = true
		return nil
	})

	ups := time.Duration(updatesPerSecond)
	ticker := time.NewTicker(time.Second / ups)
	defer ticker.Stop()

	dtt := core.NewDeltaTimeTracker()
	for range ticker.C {
		if a.exit {
			return nil
		}
		err := a.Update(dtt.Step())
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) Close() {
	for _, layer := range a.layerStack {
		layer.OnDetach()
	}
}

func NewApp(config Config) *App {
	return &App{
		config:     config,
		layerStack: make([]Layer, 0),
		exit:       false,
	}
}
