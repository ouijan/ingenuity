package cli

//go:generate echo "Run go generate ./pkg/cli to see this message"

import (
	"fmt"
	"path"
)

func generateTiled(context Context) error {
	tiledProject := path.Join(context.ProjectDir, context.Config.TiledProject)
	tiledOutput := path.Join(context.ProjectDir, context.Config.TiledTypesOutput)
	// sourcePath := path.Join(context.ProjectDir, context.Config.sourcePath)

	project, err := readTiledProject(tiledProject)
	if err != nil {
		return err
	}

	err = generateTiledCode(project, tiledOutput)
	if err != nil {
		return err
	}
	// parseSourceCode(sourcePath)
	// addTiledCommands(tiledProject)
	// addTiledProperties(tiledProject)
	// addTiledEnums(tiledProject)
	return nil
}

func addTiledCommands(tiledProject string) error {
	fmt.Printf("Adding Commands to %s\n", tiledProject)
	// TODO: Add Commands
	return nil
}

func addTiledEnums(tiledProject string) error {
	fmt.Printf("Adding Enums to %s\n", tiledProject)
	// TODO: Add Commands
	return nil
}

func addTiledProperties(tiledProject string) error {
	fmt.Printf("Adding Properties to %s\n", tiledProject)
	// TODO: Add Commands
	return nil
}
