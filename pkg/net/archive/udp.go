package networking

import (
	"errors"
	"fmt"
	"net"

	"github.com/ouijan/ingenuity/pkg/core"
)

// ---------- Message ----------

type UDPMessage struct {
	Addr    *net.UDPAddr
	Payload []byte
}

// ---------- Server ----------

type UDPServer struct {
	udpAddr string
	conn    *net.UDPConn
	MsgCh   chan UDPMessage
	quitCh  chan struct{}
	Peers   map[string]*net.UDPAddr
}

func (u *UDPServer) Start() error {
	addr, err := net.ResolveUDPAddr("udp", u.udpAddr)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	defer conn.Close()
	u.conn = conn
	core.Log.Info(fmt.Sprintf("UDP Server is listening on %s", addr))

	go u.readLoop(conn)

	<-u.quitCh
	close(u.MsgCh)
	return nil
}

func (u *UDPServer) readLoop(conn *net.UDPConn) {
	buff := make([]byte, 2048)
	for {
		n, addr, err := conn.ReadFromUDP(buff)
		if err != nil {
			fmt.Println("UDP Error:", err)
			continue
		}
		u.Peers[addr.String()] = addr
		u.MsgCh <- UDPMessage{
			Addr:    addr,
			Payload: buff[:n],
		}
	}
}

func (u *UDPServer) Write(addr *net.UDPAddr, payload []byte) error {
	_, err := u.conn.WriteToUDP(payload, addr)
	return err
}

func (u *UDPServer) Close() {
	close(u.quitCh)
}

func NewUDPServer(listenAddr string) *UDPServer {
	return &UDPServer{
		udpAddr: listenAddr,
		quitCh:  make(chan struct{}),
		MsgCh:   make(chan UDPMessage),
		Peers:   make(map[string]*net.UDPAddr),
	}
}

// ---------- Client ----------

type UDPClient struct {
	serverAddr string
	conn       *net.UDPConn
	MsgCh      chan UDPMessage
}

func (u *UDPClient) Connect() error {
	addr, err := net.ResolveUDPAddr("udp", u.serverAddr)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	u.conn = conn
	core.Log.Info(fmt.Sprintf("UDP Connection Accepted: %s", conn.RemoteAddr()))
	return nil
}

func (u *UDPClient) Listen() {
	buff := make([]byte, 2048)
	for {
		n, addr, err := u.conn.ReadFromUDP(buff)
		if errors.Is(err, net.ErrClosed) {
			continue
		}
		if err != nil {
			core.Log.Info(fmt.Sprintf("UDP Error: %s", err))
			continue
		}
		u.MsgCh <- UDPMessage{
			Addr:    addr,
			Payload: buff[:n],
		}

	}
}

func (u *UDPClient) Write(payload []byte) error {
	_, err := u.conn.Write(payload)
	return err
}

func (u *UDPClient) Close() {
	u.conn.Close()
}

func NewUDPClient(serverAddr string) *UDPClient {
	return &UDPClient{
		serverAddr: serverAddr,
		MsgCh:      make(chan UDPMessage),
	}
}
