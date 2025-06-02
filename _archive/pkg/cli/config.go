package cli

import (
	"encoding/json"
	"os"
	"path"
	"strings"
)

type Config struct {
	Name             string `json:"name"`
	TiledProject     string `json:"tiledProject"`
	TiledTypesOutput string `json:"tiledTypesOutput"`
	Src              string `json:"src"`
}

func readConfig(configPath string) (Config, error) {
	config := Config{
		Name:         "NewIngenuityProject",
		TiledProject: "tiled.tmx",
		Src:          "src",
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	decoder := json.NewDecoder(strings.NewReader(string(data)))
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	if config.TiledTypesOutput == "" {
		config.TiledTypesOutput = path.Join(config.Src, "types")
	}

	return config, nil
}

type Context struct {
	CurrentDir string
	ProjectDir string
	Config     Config
}

func buildContext(configPath string) (Context, error) {
	context := Context{}
	context.ProjectDir = path.Dir(configPath)

	currentDir, err := os.Getwd()
	if err != nil {
		return context, err
	}
	context.CurrentDir = currentDir

	config, err := readConfig(configPath)
	if err != nil {
		return Context{}, err
	}

	context.Config = config
	return context, nil
}
