package networking

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/ouijan/ingenuity/pkg/core"
)

type TCPMessage struct {
	Conn    *net.TCPConn
	Payload []byte
}

// ---------- Server ----------

type TCPServer struct {
	listenAddr     string
	listener       *net.TCPListener
	quitCh         chan struct{}
	MsgCh          chan TCPMessage
	ConnAcceptedCh chan *net.TCPConn
	ConnClosedCh   chan *net.TCPConn
	Peers          map[string]*net.TCPConn
}

func (t *TCPServer) Start() error {
	addr, err := net.ResolveTCPAddr("tcp", t.listenAddr)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	t.listener = listener
	core.Log.Info(fmt.Sprintf("TCP Server is listening on %s", addr))

	go t.acceptLoop()

	<-t.quitCh
	close(t.MsgCh)
	close(t.ConnAcceptedCh)
	close(t.ConnClosedCh)
	return nil
}

func (t *TCPServer) acceptLoop() {
	for {
		conn, err := t.listener.AcceptTCP()
		if err != nil {
			core.Log.Warn(fmt.Sprintf("TCP Error: %s", err))
			continue
		}
		t.ConnAcceptedCh <- conn
		t.Peers[conn.RemoteAddr().String()] = conn
		go t.readLoop(conn)
	}
}

func (t *TCPServer) readLoop(conn *net.TCPConn) {
	defer conn.Close()
	buff := make([]byte, 2048)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				t.ConnClosedCh <- conn
				return
			}

			core.Log.Warn(fmt.Sprintf("TCP Error:%s %s", conn.RemoteAddr(), err))
			continue
		}
		t.MsgCh <- TCPMessage{
			Conn:    conn,
			Payload: buff[:n],
		}
	}
}

func (t *TCPServer) Close() {
	close(t.quitCh)
}

func (t *TCPServer) Write(conn *net.TCPConn, payload []byte) error {
	_, err := conn.Write(payload)
	return err
}

func NewTCPServer(listenAddr string, msgBuffer int) *TCPServer {
	return &TCPServer{
		listenAddr:     listenAddr,
		quitCh:         make(chan struct{}),
		MsgCh:          make(chan TCPMessage, msgBuffer),
		ConnAcceptedCh: make(chan *net.TCPConn),
		ConnClosedCh:   make(chan *net.TCPConn),
		Peers:          make(map[string]*net.TCPConn),
	}
}

// ---------- Client ----------

type TCPClient struct {
	serverAddr string
	conn       *net.TCPConn
	MsgCh      chan TCPMessage
}

func (t *TCPClient) Connect() error {
	addr, err := net.ResolveTCPAddr("tcp", t.serverAddr)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}
	t.conn = conn
	core.Log.Info(fmt.Sprintf("TCP Connection Accepted: %s", conn.RemoteAddr()))
	return nil
}

func (t *TCPClient) Listen() {
	buff := make([]byte, 2048)
	for {
		n, err := t.conn.Read(buff)
		if err == io.EOF {
			core.Log.Info(fmt.Sprintf("TCP Connection Closed: %s", t.conn.RemoteAddr()))
			return
		}
		if errors.Is(err, net.ErrClosed) {
			continue
		}
		if err != nil {
			core.Log.Info(fmt.Sprintf("TCP Error: %s %s", t.conn.RemoteAddr(), err))
			continue
		}
		t.MsgCh <- TCPMessage{
			Conn:    t.conn,
			Payload: buff[:n],
		}
	}
}

func (t *TCPClient) Write(data []byte) error {
	_, err := t.conn.Write(data)
	return err
}

func (t *TCPClient) Close() {
	t.conn.Close()
}

func NewTCPClient(serverAddr string) *TCPClient {
	return &TCPClient{
		serverAddr: serverAddr,
		MsgCh:      make(chan TCPMessage, 10),
	}
}
