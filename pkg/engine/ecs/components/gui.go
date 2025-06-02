package components

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ouijan/ingenuity/pkg/engine/net"
	"github.com/ouijan/ingenuity/pkg/engine/utils"
)

/**
 * Text Component
 */

const TextComponentId = "comp_Text"

type Text struct {
	Content  string
	FontSize int32
	Colour   rl.Color
}

// ApplyDelta implements net.NetSyncable.
func (t *Text) ApplyDelta(delta net.SyncVars) {
	t.Content, _ = delta.GetString("content")
	t.FontSize, _ = delta.GetInt32("fontSize")
	// t.Colour, _ = delta.GetColor("colour")
}

// GetDelta implements net.NetSyncable.
func (t *Text) GetDelta() net.SyncVars {
	delta := net.NewSyncVars()
	delta.Set("content", t.Content)
	delta.Set("fontSize", t.FontSize)
	// delta.Set("colour", t.Colour)
	return delta
}

// GetId implements net.NetSyncable.
func (t *Text) GetId() utils.HashId {
	return utils.Hash(TextComponentId)
}

var _ net.Syncable = (*Text)(nil)
