package systems

import (
	ark "github.com/mlange-42/ark/ecs"
	"google.golang.org/protobuf/proto"

	"github.com/ouijan/ingenuity/pkg/engine/ecs"
	"github.com/ouijan/ingenuity/pkg/engine/ecs/components"
	"github.com/ouijan/ingenuity/pkg/engine/net"
	"github.com/ouijan/ingenuity/pkg/engine/net/packet"
	"github.com/ouijan/ingenuity/pkg/engine/utils"
)

const CommandFrameExpiry = 10

type ServerNetSync struct {
	commandFrame uint64
	server       ark.Resource[net.Server]
	netEnts      *ark.Filter1[components.NetworkedEntity]
	netComps     *ark.Map3[components.Metadata, components.Transform2D, components.Text]
}

// OnCreate implements ecs.System.
func (s *ServerNetSync) OnCreate(ea *ecs.EntityAdmin) {
	s.server = ark.NewResource[net.Server](&ea.World)
	s.netEnts = ark.NewFilter1[components.NetworkedEntity](&ea.World)
	s.netComps = ark.NewMap3[components.Metadata, components.Transform2D, components.Text](
		&ea.World,
	)
}

// OnDestroy implements ecs.System.
func (s *ServerNetSync) OnDestroy() {
	// Clean up the system resources.
}

// Update implements ecs.System.
func (s *ServerNetSync) Update(dt float32) error {
	s.commandFrame++

	netEnts := s.netEnts.Query()
	for netEnts.Next() {
		net := netEnts.Get()
		ent := netEnts.Entity()
		s.updateEntity(ent, net, s.commandFrame)
	}

	// Update the system logic.

	// This is where you would handle network synchronization.

	return nil
}

func (s *ServerNetSync) updateEntity(
	ent ark.Entity,
	ne *components.NetworkedEntity,
	commandFrame uint64,
) {
	meta, trans, text := s.netComps.Get(ent)
	comps := []any{meta, trans, text}
	syncables := net.GetSyncables(comps)

	delta := ne.SDM.GetDelta(syncables, commandFrame)

	ne.SM.AddDelta(delta)
	ne.SM.BuildPackets(commandFrame, delta)

	aks := []struct{ clientId, frame uint64 }{} // TODO: Get from network
	for _, ak := range aks {
		ne.SM.AknowledgeCommandFrame(ak.clientId, ak.frame)
	}
	ne.SM.Flush(commandFrame - CommandFrameExpiry)

	// TODO: Remove this, it sends the same frame packet to all clients
	sync := ne.SM.GetFramePacket(commandFrame, delta)
	sync.NetworkEntityId = utils.Pointer(uint64(ne.Id))

	if len(sync.Deltas) > 0 {
		p := packet.NewSyncPacket(sync)
		data, err := proto.Marshal(p)
		if err == nil {
			s.server.Get().Broadcast(data)
		}
	}

	// TODO: Use Notify iterator to get packets to send to clients
	// for notify := range ne.SM.Notify {
	// 	s.sa.server.Broadcast()
	// }

	ne.SDM.ApplyDelta(delta)
}

var _ ecs.System = (*ServerNetSync)(nil)
