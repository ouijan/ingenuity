package components

import "github.com/ouijan/ingenuity/pkg/core/net"

type NetworkedEntity struct {
	Id      uint64
	OwnerId int32
	SM      *net.SyncManager
	SDM     *net.SyncDeltaManager
}

func NewNetworkedEntity(id uint64) *NetworkedEntity {
	return &NetworkedEntity{
		Id:  id,
		SM:  net.NewSyncManager(),
		SDM: net.NewSyncDeltaManager(),
	}
}
