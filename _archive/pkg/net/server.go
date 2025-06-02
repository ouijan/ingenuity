package net

import (
	"fmt"
	"io"
	"net"

	"github.com/ouijan/ingenuity/pkg/core"
)

type server struct {
	addr           Addr
	peers          map[string]*net.TCPConn
	msgCh          chan Message
	connectedCh    chan Addr
	disconnectedCh chan Addr
	quitCh         chan struct{}
}

// OnConnect implements Server.
func (s *server) OnConnect() chan Addr {
	return s.connectedCh
}

// OnDisconnect implements Server.
func (s *server) OnDisconnect() chan Addr {
	return s.disconnectedCh
}

// Close implements Server.
func (s *server) Close() {
	close(s.quitCh)
}

// Listen implements Server.
func (s *server) Listen() error {
	addr, err := net.ResolveTCPAddr("tcp", s.addr.String())
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	core.Log.Info(fmt.Sprintf("Server is listening on %s", s.addr.String()))

	go s.acceptLoop(listener)

	<-s.quitCh
	close(s.msgCh)
	return nil
}

func (s *server) acceptLoop(listener *net.TCPListener) {
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			core.Log.Warn(fmt.Sprintf("Error: %s", err))
			continue
		}
		addr := NewAddr(conn.RemoteAddr().String())
		s.connectedCh <- addr
		s.peers[addr.String()] = conn
		go s.readLoop(addr, conn)
	}
}

func (s *server) readLoop(addr Addr, conn *net.TCPConn) {
	defer conn.Close()
	buff := make([]byte, 2048)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				s.disconnectedCh <- addr
				return
			}
			core.Log.Warn(fmt.Sprintf("Error:%s %s", conn.RemoteAddr(), err))
			continue
		}
		s.msgCh <- &message{
			addr:    addr,
			payload: buff[:n],
		}
	}
}

// Read implements Server.
func (s *server) Read() chan Message {
	return s.msgCh
}

// Write implements Server.
func (s *server) Write(addr Addr, payload []byte) error {
	conn, ok := s.peers[addr.String()]
	if !ok {
		return fmt.Errorf("peer not found")
	}
	_, err := conn.Write(payload)
	return err
}

var _ Server = (*server)(nil)

func NewServer(addr string, maxConnections int) Server {
	return &server{
		addr:           NewAddr(addr),
		peers:          make(map[string]*net.TCPConn),
		msgCh:          make(chan Message),
		connectedCh:    make(chan Addr, maxConnections),
		disconnectedCh: make(chan Addr, maxConnections),
		quitCh:         make(chan struct{}),
	}
}
