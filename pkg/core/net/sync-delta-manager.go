package net

import (
	"slices"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/ouijan/ingenuity/pkg/core/log"
	"github.com/ouijan/ingenuity/pkg/core/net/packet"
	"github.com/ouijan/ingenuity/pkg/core/utils"
)

type SyncDeltaManager struct {
	commandFrameInternal uint64
	instances            []*SyncInstance
}

func NewSyncDeltaManager() *SyncDeltaManager {
	return &SyncDeltaManager{
		commandFrameInternal: 0,
		instances:            []*SyncInstance{},
	}
}

func (sdm *SyncDeltaManager) FindInstance(id utils.HashId) *SyncInstance {
	for _, instance := range sdm.instances {
		if instance.instanceId == id {
			return instance
		}
	}
	return nil
}

func (sdm *SyncDeltaManager) GetDelta(
	comps []Syncable,
	commandFrame uint64,
) *packet.Sync_FrameDelta {
	changes := []*packet.Sync_InstanceDelta{}

	for _, comp := range comps {
		change := sdm.buildDeltaChange(comp)
		if change == nil {
			continue
		}
		isCreated := change.Created != nil && *change.Created
		isDestroyed := change.Destroyed != nil && *change.Destroyed

		if isCreated || isDestroyed || change.Changes != nil {
			changes = append(changes, change)
		}
	}

	for _, instance := range sdm.instances {
		var comp Syncable
		for _, c := range comps {
			if c.GetId() == instance.instanceId {
				comp = c
				break
			}
		}
		if comp == nil {
			change := &packet.Sync_InstanceDelta{
				InstanceId: utils.Pointer(uint64(instance.instanceId)),
				Destroyed:  utils.Pointer(true),
			}
			changes = append(changes, change)
		}
	}

	return &packet.Sync_FrameDelta{
		CommandFrame: &commandFrame,
		Changes:      changes,
	}
}

func (sdm *SyncDeltaManager) buildDeltaChange(comp Syncable) *packet.Sync_InstanceDelta {
	vars := comp.GetDelta()

	changes, err := structpb.NewStruct(vars.Data())
	if err != nil {
		log.Error("Failed to create changes struct: %v", err)
		return nil
	}

	iDelta := &packet.Sync_InstanceDelta{
		InstanceId: utils.Pointer(uint64(comp.GetId())),
		Changes:    changes,
	}

	instance := sdm.FindInstance(comp.GetId())
	if instance == nil {
		iDelta.Created = utils.Pointer(true)
		return iDelta
	}

	if vars.IsEq(instance.instanceVars) {
		return nil
	}

	return iDelta
}

func (sdm *SyncDeltaManager) ApplyDelta(delta *packet.Sync_FrameDelta) {
	if sdm.commandFrameInternal >= *delta.CommandFrame {
		return
	}

	for _, change := range delta.Changes {
		if change.Destroyed != nil && *change.Destroyed {
			sdm.instances = slices.DeleteFunc(sdm.instances, func(i *SyncInstance) bool {
				return i.instanceId == utils.HashId(*change.InstanceId)
			})
			continue
		}

		if change.Created != nil && *change.Created {
			sdm.instances = append(sdm.instances, &SyncInstance{
				instanceId:   utils.HashId(*change.InstanceId),
				instanceVars: SyncVarsFromMap(change.Changes.AsMap()),
			})
			continue
		}

		instance := sdm.FindInstance(*change.InstanceId)
		instance.instanceVars = SyncVarsFromMap(change.Changes.AsMap())
	}

	sdm.commandFrameInternal = *delta.CommandFrame
}
