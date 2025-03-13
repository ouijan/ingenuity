package engine

import "github.com/ouijan/ingenuity/pkg/resources"

type TransformComponent struct {
	X, Y float64
}

type SpriteRendererComponent struct {
	SpriteSheet *resources.SpriteSheet
	SpriteIndex int
}

type (
	ColliderCategory       uint
	CollisionMask          uint
	BoxCollider2DComponent struct {
		T, B, L, R   float64
		Category     ColliderCategory
		CategoryMask CollisionMask
	}
)

type rigidBodyType uint

const (
	RB_Static rigidBodyType = iota
	RB_Dynamic
	RB_Kinematic
)

type RigidBody2DComponent struct {
	Type   rigidBodyType
	Mass   float64
	Vx, Vy float64
}
