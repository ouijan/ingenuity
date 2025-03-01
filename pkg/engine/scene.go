package engine

type IScene interface {
	Load()
	OnEnter(world *IWorld)
	OnExit(world *IWorld)
}

type SceneManager struct {
	active    IScene
	nextScene IScene
}

func (sm *SceneManager) Active() IScene {
	return sm.active
}

func (sm *SceneManager) SetNext(scene IScene) {
	sm.nextScene = scene
}

func (sm *SceneManager) Update() {
	if sm.active == sm.nextScene {
		return
	}
	if sm.nextScene != nil {
		sm.setActive(sm.nextScene)
		sm.nextScene = nil
	}
}

func (sm *SceneManager) setActive(scene IScene) {
	scene.Load()
	if sm.active != nil {
		sm.active.OnExit(World)
	}
	World = NewWorld() // TODO: I don't like that this single is being overwritten in the global scope
	scene.OnEnter(World)
	sm.active = scene
}

func NewSceneManager() *SceneManager {
	return &SceneManager{}
}

var Scene = NewSceneManager()
