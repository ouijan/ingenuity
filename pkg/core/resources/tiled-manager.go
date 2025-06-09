package resources

import (
	"os"

	"github.com/lafriks/go-tiled"
)

type tiledMapManager struct {
	root string
}

func newTiledMapManager(root string) *tiledMapManager {
	return &tiledMapManager{root}
}

func (m *tiledMapManager) Load(name string) (*tiled.Map, error) {
	fileSystem := os.DirFS(m.root)
	tilemap, err := tiled.LoadFile(name, tiled.WithFileSystem(fileSystem))
	return tilemap, err
}
