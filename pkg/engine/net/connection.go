package net

import (
	"context"
	"sync"
	"time"

	"github.com/pion/dtls/v3"

	"github.com/ouijan/ingenuity/pkg/engine/log"
)

type Conn struct {
	*dtls.Conn
	buffer *PacketBuffer
}

func (c *Conn) Handshake() error {
	timeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	err := c.HandshakeContext(ctx)
	cancel()

	return err
}

func (c *Conn) Listen(pb *PacketBuffer) error {
	buffer := make([]byte, pb.BufferSize)
	for {
		n, err := c.Read(buffer)
		if err != nil {
			return err
		}
		pb.Write(c.RemoteAddr(), buffer[:n])
	}
}

func NewConnection(conn *dtls.Conn) *Conn {
	return &Conn{
		Conn:   conn,
		buffer: NewPacketBuffer(8192),
	}
}

// --

type ConnManager struct {
	conns map[string]*Conn
	lock  sync.RWMutex
}

func (c *ConnManager) Register(conn *Conn) {
	log.Info("Connected to %s", conn.RemoteAddr())

	c.lock.Lock()
	defer c.lock.Unlock()

	c.conns[conn.RemoteAddr().String()] = conn
}

func (c *ConnManager) Unregister(conn *Conn) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.conns, conn.RemoteAddr().String())
	err := conn.Close()

	if err != nil {
		log.Error("Failed to disconnect", conn.RemoteAddr(), err)
	} else {
		log.Info("Disconnected ", conn.RemoteAddr())
	}
}

func (c *ConnManager) Connections(yield func(*Conn) bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for _, s := range c.conns {
		if !yield(s) {
			return
		}
	}
}

func NewConnectionManager() *ConnManager {
	return &ConnManager{
		conns: make(map[string]*Conn),
	}
}
