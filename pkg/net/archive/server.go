package networking

import (
	"fmt"

	"github.com/ouijan/ingenuity/pkg/core"
)

type Server struct {
	quitCh chan struct{}
	TCP    *TCPServer
	UDP    *UDPServer
}

func (s *Server) Start() {
	go func() {
		for conn := range s.TCP.ConnAcceptedCh {
			core.Log.Info(fmt.Sprintf("TCP Connection Accepted: %s", conn.RemoteAddr()))
		}
	}()
	go func() {
		for conn := range s.TCP.ConnClosedCh {
			core.Log.Info(fmt.Sprintf("TCP Connection Closed: %s", conn.RemoteAddr()))
		}
	}()

	go s.TCP.Start()
	defer s.TCP.Close()

	go s.UDP.Start()
	defer s.UDP.Close()

	<-s.quitCh
}

func (s *Server) Close() {
	close(s.quitCh)
}

func NewServer(port int) *Server {
	host := fmt.Sprintf("localhost:%d", port)
	tcpBuffer := 10

	return &Server{
		TCP: NewTCPServer(host, tcpBuffer),
		UDP: NewUDPServer(host),
	}
}
