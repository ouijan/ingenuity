package net

import (
	"google.golang.org/protobuf/proto"

	"github.com/ouijan/ingenuity/pkg/core/net/packet"
	"github.com/ouijan/ingenuity/pkg/core/utils"
)

// Maintained until client disconnects
type SyncGhost struct {
	clientId            uint64
	commandFrameLastAck uint64
	outstandingPackets  []*packet.Sync
}

func SyncMarshal(packet *packet.Sync) ([]byte, error) {
	out, err := proto.Marshal(packet)
	return out, err
}

// ----- ECS

type Syncable interface {
	GetId() utils.HashId
	ApplyDelta(delta SyncVars)
	GetDelta() SyncVars
}

func GetSyncables(comps []interface{}) []Syncable {
	var syncables []Syncable
	for _, comp := range comps {
		if syncable, ok := comp.(Syncable); ok {
			syncables = append(syncables, syncable)
		}
	}
	return syncables
}

// --- CLient side Pseudocode
func ReplicatePacket(sdm *SyncDeltaManager, packet *packet.Sync) {
	// TODO: Send acknowledgement
	if sdm.commandFrameInternal >= *packet.CommandFrameEnd {
		return
	}

	if *packet.IsLocal {
		// TODO: Rollback
	}

	// Assume sorted sequentially
	for _, delta := range packet.Deltas {
		replicateDelta(sdm, delta)
		sdm.ApplyDelta(delta)
	}

	if *packet.IsLocal {
		// TODO: Simulate
	}
}

func replicateDelta(sdm *SyncDeltaManager, delta *packet.Sync_FrameDelta) {
	if sdm.commandFrameInternal >= *delta.CommandFrame {
		return
	}

	for _, change := range delta.Changes {
		if *change.Destroyed {
			// TODO: que component for removal
			continue
		}

		if *change.Created {
			// TODO: create a new instance
			// needs to resolve a create functon
			continue
		}

		// var comp NetSyncable // TODO: Resolve component
		// if comp != nil {
		// 	comp.ApplyDelta(change.varChanges)
		// }
	}
}
