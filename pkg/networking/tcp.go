package networking

import (
	"fmt"
	"io"
	"net"
)

type TCPMessage struct {
	Conn    net.Conn
	Payload []byte
}

// ---------- Server ----------
type TCPServer struct {
	listenAddr string
	listener   net.Listener
	quitCh     chan struct{}
	MsgCh      chan TCPMessage
	Peers      map[string]net.Conn
}

func (t *TCPServer) Start() error {
	listener, err := net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	t.listener = listener
	fmt.Println("Server is listening on port 8080")

	go t.acceptLoop()

	<-t.quitCh
	t.Peers = make(map[string]net.Conn)
	close(t.MsgCh)
	return nil
}

func (t *TCPServer) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		logCon(conn, "Accepted connection")
		t.Peers[conn.RemoteAddr().String()] = conn
		go t.readLoop(conn)
	}
}

func (t *TCPServer) readLoop(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte, 2048)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				logCon(conn, "Connection closed")
				delete(t.Peers, conn.RemoteAddr().String())
				return
			}
			logCon(conn, "Error: %s", err)
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

func NewTCPServer(listenAddr string, msgBuffer int) *TCPServer {
	return &TCPServer{
		listenAddr: listenAddr,
		quitCh:     make(chan struct{}),
		MsgCh:      make(chan TCPMessage, msgBuffer),
		Peers:      make(map[string]net.Conn),
	}
}

// ---------- Client ----------

type TCPClient struct {
	serverAddr string
	conn       net.Conn
	MsgCh      chan TCPMessage
}

func (t *TCPClient) Connect() error {
	conn, err := net.Dial("tcp", t.serverAddr)
	if err != nil {
		return err
	}
	// defer conn.Close()
	t.conn = conn

	logCon(conn, "Connected to server")
	return nil
}

func (t *TCPClient) ReadLoop() {
	buff := make([]byte, 2048)
	for {
		n, err := t.conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				logCon(t.conn, "Connection closed")
				return
			}
			logCon(t.conn, "Error: %s", err)
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

// ---------- Helpers ----------

func logCon(conn net.Conn, msg string, args ...interface{}) {
	fmt.Printf("[%s] %s\n", conn.RemoteAddr(), fmt.Sprintf(msg, args...))
}
