package net

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/pion/dtls/v3"
	"github.com/pion/dtls/v3/pkg/crypto/selfsign"

	"github.com/ouijan/ingenuity/pkg/core/log"
)

type Server struct {
	port     int16
	Buffer   *PacketBuffer
	quitCh   chan struct{}
	readyCh  chan struct{}
	cm       *ConnManager
	listener net.Listener
}

func (s *Server) Start() error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	certificate, err := selfsign.GenerateSelfSigned()
	if err != nil {
		return err
	}

	config := &dtls.Config{
		Certificates:         []tls.Certificate{certificate},
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
	}

	listener, err := dtls.Listen("udp", addr, config)
	if err != nil {
		return err
	}
	s.listener = listener
	return nil
}

func (s *Server) Listen() error {
	defer s.listener.Close()

	log.Info("Listening on port %d", s.port)
	go s.acceptConns()

	<-s.quitCh
	s.Buffer.Close()
	return nil
}

func (s *Server) acceptConns() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Error("Failed to accept connection: %v", err)
			return
		}
		go s.handleConn(conn)

	}
}

func (s *Server) handleConn(netConn net.Conn) {
	dtlsConn, ok := netConn.(*dtls.Conn)
	if !ok {
		log.Error("Connection is not a DTLS connection")
		return
	}

	conn := NewConnection(dtlsConn)
	if err := conn.Handshake(); err != nil {
		log.Error("Failed to perform DTLS handshake: %v", err)
		return
	}

	s.cm.Register(conn)
	defer s.cm.Unregister(conn)

	err := conn.Listen(conn.buffer)
	if err != nil {
		log.Error("Failed to listen on connection: %v", err)
		return
	}
}

func (s *Server) Broadcast(payload []byte) {
	// log.Info("-> Broadcasting message: %s", string(payload))
	for conn := range s.cm.Connections {
		_, err := conn.Write(payload)
		if err != nil {
			log.Warn("Failed to write message to %s: %v", conn.RemoteAddr(), err)
		}
	}
}

func (u *Server) Close() {
	close(u.quitCh)
}

func NewServer(port int16) *Server {
	return &Server{
		port:   port,
		Buffer: NewPacketBuffer(8192),
		quitCh: make(chan struct{}),
		cm:     NewConnectionManager(),
	}
}

func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}
