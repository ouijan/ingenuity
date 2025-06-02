package net

import (
	"slices"

	"github.com/ouijan/ingenuity/pkg/engine/net/packet"
	"github.com/ouijan/ingenuity/pkg/engine/utils"
)

type SyncManager struct {
	deltas  []*packet.Sync_FrameDelta
	ghosts  []SyncGhost
	packets []*packet.Sync
}

func NewSyncManager() *SyncManager {
	return &SyncManager{
		deltas:  []*packet.Sync_FrameDelta{},
		ghosts:  []SyncGhost{},
		packets: []*packet.Sync{},
	}
}

func (sm *SyncManager) FindGhost(clientId uint64) *SyncGhost {
	for _, ghost := range sm.ghosts {
		if ghost.clientId == clientId {
			return &ghost
		}
	}
	return nil
}

func (sm *SyncManager) AddDelta(delta *packet.Sync_FrameDelta) {
	sm.deltas = append(sm.deltas, delta)
}

func (sm *SyncManager) Flush(commandFrameExpiry uint64) {
	lastAck := sm.getOldestAcknowledgedFrame()

	sm.packets = slices.DeleteFunc(sm.packets, func(packet *packet.Sync) bool {
		return *packet.CommandFrameEnd <= lastAck || *packet.CommandFrameEnd <= commandFrameExpiry
	})

	for _, ghost := range sm.ghosts {
		ghost.outstandingPackets = slices.DeleteFunc(
			ghost.outstandingPackets,
			func(i *packet.Sync) bool {
				return i == nil || *i.CommandFrameEnd <= lastAck ||
					*i.CommandFrameEnd <= ghost.commandFrameLastAck
			},
		)
	}

	sm.deltas = slices.DeleteFunc(sm.deltas, func(delta *packet.Sync_FrameDelta) bool {
		return *delta.CommandFrame <= lastAck
	})
}

func (sm *SyncManager) GetFramePacket(
	commandFrame uint64,
	delta *packet.Sync_FrameDelta,
) *packet.Sync {
	return &packet.Sync{
		IsLocal:           utils.Pointer(false),
		CommandFrameStart: &commandFrame,
		CommandFrameEnd:   &commandFrame,
		Deltas:            []*packet.Sync_FrameDelta{delta},
	}
}

func (sm *SyncManager) BuildPackets(commandFrame uint64, frameDelta *packet.Sync_FrameDelta) {
	packet := sm.GetFramePacket(commandFrame, frameDelta)
	sm.packets = append(sm.packets, packet)

	// TODO: We can broadcast for now, but will need to tailor packets for each client (ghost)
	// for _, ghost := range syncer.syncManager.ghosts {
	// 	packet := &NetSyncPacket{
	// 		isLocal:           false, // check if ghost.clientId owns this entity
	// 		commandFrameStart: commandFrame,
	// 		commandFrameEnd:   commandFrame,
	// 		deltas:            []NetSyncDelta{frameDelta},
	// 	}
	// 	ghost.outstandingPackets = append(ghost.outstandingPackets, packet)
	// 	syncer.syncManager.packets = append(syncer.syncManager.packets, packet)
	// }
}

func (sm *SyncManager) AknowledgeCommandFrame(clientId, commandFrame uint64) {
	ghost := sm.FindGhost(clientId)
	if ghost == nil {
		return
	}
	ghost.commandFrameLastAck = commandFrame
}

func (sm *SyncManager) getOldestAcknowledgedFrame() uint64 {
	oldest := uint64(0)
	for _, ghost := range sm.ghosts {
		if oldest == 0 || ghost.commandFrameLastAck < oldest {
			oldest = ghost.commandFrameLastAck
		}
	}
	return oldest
}

type SyncNotify struct {
	Ghost  *SyncGhost
	Packet *packet.Sync
}

func (sm *SyncManager) Notify(yield func(SyncNotify) bool) {
	for _, ghost := range sm.ghosts {
		for _, packet := range ghost.outstandingPackets {
			notify := SyncNotify{Ghost: &ghost, Packet: packet}
			if yield(notify) {
				continue
			}
			return
		}
	}
}
