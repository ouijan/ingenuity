package systems

import (
	"fmt"
	"time"

	ark "github.com/mlange-42/ark/ecs"

	"github.com/ouijan/ingenuity/pkg/core/ecs"
	"github.com/ouijan/ingenuity/pkg/core/ecs/components"
	"github.com/ouijan/ingenuity/pkg/core/log"
	"github.com/ouijan/ingenuity/pkg/core/net"
	"github.com/ouijan/ingenuity/pkg/core/net/packet"
	"github.com/ouijan/ingenuity/pkg/core/utils"
)

type NetSync struct {
	client   ark.Resource[net.Client]
	netEnts  *ark.Filter1[components.NetworkedEntity]
	netComps *ark.Map3[components.Metadata, components.Transform2D, components.Text]
}

// OnCreate implements ecs.System.
func (s *NetSync) OnCreate(ea *ecs.EntityAdmin) {
	s.client = ark.NewResource[net.Client](&ea.World)

	s.netEnts = ark.NewFilter1[components.NetworkedEntity](&ea.World)
	s.netComps = ark.NewMap3[components.Metadata, components.Transform2D, components.Text](
		&ea.World,
	)
}

// OnDestroy implements ecs.System.
func (s *NetSync) OnDestroy() {
	// Clean up the system resources.
}

// Update implements ecs.System.
func (s *NetSync) Update(dt float32) error {
	client := s.client.Get()
	utils.Assert(client != nil)

	for {
		select {
		case _ = <-client.Buffer.AckCh:
			// log.Debug("Received ack packet: %s", ack.String())
			continue
		case _ = <-client.Buffer.MsgCh:
			// log.Debug("Received message packet: %s", msg.String())
			continue
		case sync := <-client.Buffer.SyncCh:
			// log.Debug("Received sync packet: %s", sync.String())
			s.handleSyncPacket(sync)
			continue
		case <-time.After(utils.ServerTickRate / 4):
			// log.Debug("Timeout waiting for sync packet, breaking the loop.")
			fmt.Println("timeout")
			return nil
		default:
			// log.Debug("No sync packets received, checking the channel again.")
			// No messages in the channel, break the loop.
			// This prevents blocking on the channel.
			// You can also use a timeout or a condition to break the loop.
			// This is useful if you want to process other events or perform other tasks.
			// log.Debug("No messages in the channel, breaking the loop.")
			return nil
		}
	}
}

func (s *NetSync) handleSyncPacket(p net.InboundPacket[*packet.Sync]) error {
	netEnts := s.netEnts.Query()
	for netEnts.Next() {
		net := netEnts.Get()
		ent := netEnts.Entity()
		if net.Id == *p.Data.NetworkEntityId {
			err := s.updateEntity(ent, net, p.Data)
			if err != nil {
				log.Error("Failed to update entity: %v", err)
				return err
			}
		}
	}

	return nil
}

func (s *NetSync) updateEntity(
	ent ark.Entity,
	ne *components.NetworkedEntity,
	sync *packet.Sync,
) error {
	meta, trans, text := s.netComps.Get(ent)
	comps := []any{meta, trans, text}
	syncables := net.GetSyncables(comps)

	for _, delta := range sync.Deltas {
		// delta.CommandFrame
		for _, change := range delta.Changes {

			if change.GetCreated() {
				log.Debug("Creating syncable %s for entity %d", change.GetInstanceId(), ne.Id)
				continue
			}

			if change.GetDestroyed() {
				log.Debug("Destroying syncable %s for entity %d", change.GetInstanceId(), ne.Id)
				continue
			}

			var targetSyncable net.Syncable
			for _, syncable := range syncables {
				if syncable.GetId() == change.GetInstanceId() {
					targetSyncable = syncable
					break
				}
			}

			if targetSyncable == nil {
				log.Warn("Syncable %s not found for entity %d", change.GetInstanceId(), ne.Id)
				continue
			}

			syncVars := net.SyncVarsFromMap(change.GetChanges().AsMap())
			log.Debug(
				"Applying delta to syncable %s for entity %d: %v",
				change.GetInstanceId(),
				ne.Id,
				syncVars,
			)
			targetSyncable.ApplyDelta(syncVars)
		}

		// ne.SM.AknowledgeCommandFrame()
	}

	return nil
}

var _ ecs.System = (*NetSync)(nil)
