package main

import (
	"github.com/ouijan/ingenuity/pkg/client"
	"github.com/ouijan/ingenuity/pkg/core/config"
)

func main() {
	cfg := config.NewConfig()
	client := client.NewClientApp(cfg)
	defer client.Close()

	client.Init()
	if err := client.Run(); err != nil {
		panic(err)
	}
}
