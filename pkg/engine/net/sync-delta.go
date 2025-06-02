package net

import (
	"github.com/ouijan/ingenuity/pkg/engine/utils"
)

// -----

// Represents a single script (component) on the entity
// Can be synchronized or unsynchronized
type SyncInstance struct {
	instanceId   utils.HashId
	instanceVars SyncVars
	// stuGraph     *interface{}
	// states []*NetSyncState
	// futureEvents []*interface{}
}

func NewSyncInstance(id utils.HashId) *SyncInstance {
	return &SyncInstance{
		instanceId:   id,
		instanceVars: NewSyncVars(),
	}
}
