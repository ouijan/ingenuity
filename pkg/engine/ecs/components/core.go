package components

import (
	"github.com/ouijan/ingenuity/pkg/engine/log"
	"github.com/ouijan/ingenuity/pkg/engine/net"
	"github.com/ouijan/ingenuity/pkg/engine/utils"
)

/**
 * Metadata Component
 */

const MetadataComponentId = "comp_Metadata"

type Metadata struct {
	Name string
	Tags []string
}

// ApplyDelta implements net.NetSyncable.
func (m *Metadata) ApplyDelta(delta net.SyncVars) {
	m.Name, _ = delta.GetString("name")
	// m.Tags, _ = delta.GetStringSlice("tags")
}

// GetDelta implements net.NetSyncable.
func (m *Metadata) GetDelta() net.SyncVars {
	delta := net.NewSyncVars()
	delta.Set("name", m.Name)
	// delta.Set("tags", m.Tags)
	return delta
}

// GetId implements net.NetSyncable.
func (m *Metadata) GetId() utils.HashId {
	return utils.Hash(MetadataComponentId)
}

var _ net.Syncable = (*Metadata)(nil)

/**
 * Transform2D Component
 */

const Transform2DComponentId = "comp_Transform2D"

type Transform2D struct {
	X, Y float32
	// Scale float32
	// Rotation float32
}

// ApplyDelta implements net.NetSyncable.
func (t *Transform2D) ApplyDelta(delta net.SyncVars) {
	t.X, _ = delta.GetFloat32("X")
	t.Y, _ = delta.GetFloat32("Y")
	log.Debug("Transform2D ApplyDelta: X=%.2f, Y=%.2f", t.X, t.Y)
}

// GetDelta implements net.NetSyncable.
func (t *Transform2D) GetDelta() net.SyncVars {
	delta := net.NewSyncVars()
	delta.Set("X", t.X)
	delta.Set("Y", t.Y)
	return delta
}

// GetId implements net.NetSyncable.
func (t *Transform2D) GetId() utils.HashId {
	return utils.Hash(Transform2DComponentId)
}

var _ net.Syncable = (*Transform2D)(nil)
