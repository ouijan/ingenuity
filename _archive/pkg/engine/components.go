package engine

import "github.com/ouijan/ingenuity/pkg/resources"

// TransformComponent

type TransformComponent struct {
	X, Y float64
}

// SpriteRendererComponent

type SpriteRendererComponent struct {
	SpriteSheet *resources.SpriteSheet
	SpriteIndex int
}

// ColliderComponent

type (
	ColliderCategory uint
	CollisionMask    uint
)

type CollisionEvent struct {
	TargetEntity Entity
	OtherEntity  Entity
}

type BoxCollider2DComponent struct {
	T, B, L, R   float64
	Category     ColliderCategory
	CategoryMask CollisionMask
	Collisions   []CollisionEvent
}

// RigidBodyComponent

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

// PlayerSpawnComponent
