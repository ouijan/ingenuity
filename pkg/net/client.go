package net

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/ouijan/ingenuity/pkg/core"
)

type client struct {
	remoteAddr Addr
	msgCh      chan Message
	conn       *net.TCPConn
}

// Close implements Client.
func (c *client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// Connect implements Client.
func (c *client) Connect() error {
	addr, err := net.ResolveTCPAddr("tcp", c.remoteAddr.String())
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}
	c.conn = conn
	core.Log.Info(fmt.Sprintf("Connection Accepted: %s", conn.RemoteAddr()))
	return nil
}

// Listen implements Client.
func (c *client) Listen() {
	buff := make([]byte, 2048)
	for {
		n, err := c.conn.Read(buff)
		if err == io.EOF {
			core.Log.Info(fmt.Sprintf("Connection Closed: %s", c.remoteAddr))
			return
		}
		if errors.Is(err, net.ErrClosed) {
			continue
		}
		if err != nil {
			core.Log.Info(fmt.Sprintf("Error: %s %s", c.remoteAddr, err))
			continue
		}
		c.msgCh <- &message{
			addr:    c.remoteAddr,
			payload: buff[:n],
		}
	}
}

// Read implements Client.
func (c *client) Read() chan Message {
	return c.msgCh
}

// Write implements Client.
func (c *client) Write(payload []byte) error {
	_, err := c.conn.Write(payload)
	return err
}

var _ Client = (*client)(nil)

func NewClient(addr string) Client {
	return &client{
		remoteAddr: NewAddr(addr),
		msgCh:      make(chan Message, 10),
	}
}
