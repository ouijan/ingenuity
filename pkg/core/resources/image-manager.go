package resources

import (
	"path"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type imageManager struct {
	root string
}

func newImageManager(root string) *imageManager {
	return &imageManager{root}
}

func (m *imageManager) Load(name string) (*rl.Image, error) {
	filePath := path.Join(m.root, name)
	img := rl.LoadImage(filePath)
	return img, nil
}
