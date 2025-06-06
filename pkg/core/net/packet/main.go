package packet

//go:generate protoc --go_out=paths=source_relative:. -I. ./packet.proto

func NewSyncPacket(data *Sync) *Packet {
	return &Packet{
		PacketData: &Packet_Sync{
			Sync: data,
		},
	}
}
