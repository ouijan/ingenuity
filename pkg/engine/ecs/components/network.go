package components

import "github.com/ouijan/ingenuity/pkg/engine/net"

type NetworkedEntity struct {
	Id      int32
	OwnerId int32
	SM      *net.SyncManager
	SDM     *net.SyncDeltaManager
}

func NewNetworkedEntity(id int32) *NetworkedEntity {
	return &NetworkedEntity{
		Id:  id,
		SM:  net.NewSyncManager(),
		SDM: net.NewSyncDeltaManager(),
	}
}
