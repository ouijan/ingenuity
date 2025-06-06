package components

import (
	"github.com/ouijan/ingenuity/pkg/core/log"
	"github.com/ouijan/ingenuity/pkg/core/net"
	"github.com/ouijan/ingenuity/pkg/core/utils"
)

/**
 * Text Component
 */

const TextComponentId = "comp_Text"

type Text struct {
	Content  string
	FontSize int32
	// Colour   rl.Color
}

// ApplyDelta implements net.NetSyncable.
func (t *Text) ApplyDelta(delta net.SyncVars) {
	t.Content, _ = delta.GetString("content")
	t.FontSize, _ = delta.GetInt32("fontSize")
	log.Debug("Text ApplyDelta: Content='%s', FontSize=%d", t.Content, t.FontSize)
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
