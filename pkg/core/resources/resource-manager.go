package resources

type ResourceManager struct {
	Root     string
	Image    *imageManager
	TiledMap *tiledMapManager
}

func NewResourceManager(root string) *ResourceManager {
	return &ResourceManager{
		Root:     root,
		Image:    newImageManager(root),
		TiledMap: newTiledMapManager(root),
	}
}
