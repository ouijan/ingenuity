package main

import (
	"github.com/ouijan/ingenuity/pkg/core/config"
	"github.com/ouijan/ingenuity/pkg/server"
)

func main() {
	cfg := config.NewConfig()
	app := server.NewServerApp(cfg)
	// defer app.Close()

	err := app.Init()
	if err != nil {
		panic(err)
	}

	// Run the application
	if err := app.Run(); err != nil {
		panic(err)
	}
}
