package net

import (
	"crypto/tls"
	"errors"
	"net"

	"github.com/ouijan/ingenuity/pkg/engine/log"
	"github.com/pion/dtls/v3"
	"github.com/pion/dtls/v3/pkg/crypto/selfsign"
)

type Client struct {
	bufferSize int
	serverAddr string
	conn       *Conn
	Buffer     *PacketBuffer
}

func (c *Client) Connect() error {
	addr, err := net.ResolveUDPAddr("udp", c.serverAddr)
	if err != nil {
		return err
	}

	certificate, err := selfsign.GenerateSelfSigned()
	if err != nil {
		return err
	}

	config := &dtls.Config{
		Certificates:         []tls.Certificate{certificate},
		InsecureSkipVerify:   true,
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
	}

	dtlsConn, err := dtls.Dial("udp", addr, config)
	if err != nil {
		return err
	}

	conn := NewConnection(dtlsConn)
	if err := conn.Handshake(); err != nil {
		log.Error("Failed to handshake with server: %v", err)
		return err
	}

	c.conn = conn
	log.Info("Connected to %s", addr)

	return nil
}

func (c *Client) Listen() error {
	if c.conn == nil {
		return errors.New("client not connected")
	}
	return c.conn.Listen(c.Buffer)
}

func (c *Client) Write(payload []byte) error {
	if c.conn == nil {
		return errors.New("client not connected")
	}
	log.Info("-> Sending message to %s: %s", c.serverAddr, string(payload))
	_, err := c.conn.Write(payload)
	return err
}

func (u *Client) Close() {
	if u.conn != nil {
		u.conn.Close()
	}
}

func NewClient(serverAddr string) *Client {
	return &Client{
		bufferSize: 16384,
		serverAddr: serverAddr,
		Buffer:     NewPacketBuffer(8192),
	}
}
