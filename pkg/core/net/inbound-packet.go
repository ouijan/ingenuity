package net

import (
	"net"

	"google.golang.org/protobuf/proto"

	"github.com/ouijan/ingenuity/pkg/core/log"
	"github.com/ouijan/ingenuity/pkg/core/net/packet"
)

type InboundPacket[T any] struct {
	Origin  net.Addr
	Payload []byte
	Packet  *packet.Packet
	Data    T
}

func (p *InboundPacket[T]) String() string {
	return p.Packet.String()
}

type PacketBuffer struct {
	BufferSize int
	AckCh      chan InboundPacket[*packet.Acknowledgement]
	MsgCh      chan InboundPacket[*packet.Message]
	SyncCh     chan InboundPacket[*packet.Sync]
}

func NewPacketBuffer(bufferSize int) *PacketBuffer {
	return &PacketBuffer{
		BufferSize: bufferSize,
		AckCh:      make(chan InboundPacket[*packet.Acknowledgement], bufferSize),
		MsgCh:      make(chan InboundPacket[*packet.Message], bufferSize),
		SyncCh:     make(chan InboundPacket[*packet.Sync], bufferSize),
	}
}

func (pc *PacketBuffer) Write(
	origin net.Addr,
	payload []byte,
) {
	p := &packet.Packet{}
	if err := proto.Unmarshal(payload, p); err != nil {
		log.Error("Failed to unmarshal packet: %s", err)
	}

	ip := InboundPacket[any]{
		Origin:  origin,
		Payload: payload,
		Packet:  p,
		Data:    nil,
	}

	switch d := p.PacketData.(type) {
	case *packet.Packet_Acknowledgement:
		pc.AckCh <- toTypedPacket(ip, d.Acknowledgement)
		return
	case *packet.Packet_Sync:
		pc.SyncCh <- toTypedPacket(ip, d.Sync)
		return
	case *packet.Packet_Message:
		pc.MsgCh <- toTypedPacket(ip, d.Message)
		return
	default:
		log.Error("Dropping packet unknown type: %T", p.PacketData)
	}
}

func (pb *PacketBuffer) Close() {
	close(pb.AckCh)
	close(pb.MsgCh)
	close(pb.SyncCh)
}

func toTypedPacket[T any](p InboundPacket[any], data T) InboundPacket[T] {
	return InboundPacket[T]{
		Origin:  p.Origin,
		Payload: p.Payload,
		Packet:  p.Packet,
		Data:    data,
	}
}
