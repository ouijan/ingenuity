package systems

import (
	"fmt"
	"time"

	ark "github.com/mlange-42/ark/ecs"

	"github.com/ouijan/ingenuity/pkg/engine/ecs"
	"github.com/ouijan/ingenuity/pkg/engine/log"
	"github.com/ouijan/ingenuity/pkg/engine/net"
	"github.com/ouijan/ingenuity/pkg/engine/net/packet"
	"github.com/ouijan/ingenuity/pkg/engine/utils"
)

type ClientNetSync struct {
	client ark.Resource[net.Client]
}

// OnCreate implements ecs.System.
func (s *ClientNetSync) OnCreate(ea *ecs.EntityAdmin) {
	s.client = ark.NewResource[net.Client](&ea.World)

	// Initialize the system with the entity admin.
}

// OnDestroy implements ecs.System.
func (s *ClientNetSync) OnDestroy() {
	// Clean up the system resources.
}

// Update implements ecs.System.
func (s *ClientNetSync) Update(dt float32) error {
	client := s.client.Get()
	utils.Assert(client != nil)

	for {
		select {
		case p := <-client.Buffer.SyncCh:
			s.handleSyncPacket(p)
			continue
		case <-time.After(utils.ServerTickRate / 4):
			fmt.Println("timeout")
			return nil
		default:

			// No messages in the channel, break the loop.
			// This prevents blocking on the channel.
			// You can also use a timeout or a condition to break the loop.
			// This is useful if you want to process other events or perform other tasks.
			// log.Debug("No messages in the channel, breaking the loop.")
			return nil
		}
	}
}

func (s *ClientNetSync) handleSyncPacket(p net.InboundPacket[*packet.Sync]) error {
	log.Debug("Received message %s", p.String())
	// Handle the sync packet.
	// This is where you would apply the changes to the local entities.
	return nil
}

var _ ecs.System = (*ClientNetSync)(nil)
