package net

type Addr interface {
	// GetHost() string
	// GetPort() int
	String() string
}

type Message interface {
	GetAddr() Addr
	GetPayload() []byte
}

type Server interface {
	Listen() error
	Write(addr Addr, payload []byte) error
	OnConnect() chan Addr
	OnDisconnect() chan Addr
	Read() chan Message
	Close()
}

type Client interface {
	Connect() error
	Listen()
	Read() chan Message
	Write(payload []byte) error
	Close()
}

// ---------- Message ----------

type message struct {
	addr    Addr
	payload []byte
}

func (m *message) GetAddr() Addr {
	return m.addr
}

func (m *message) GetPayload() []byte {
	return m.payload
}

var _ Message = (*message)(nil)

// ---------- Addr ----------

type addr struct {
	host string
}

func (a *addr) String() string {
	return a.host
}

var _ Addr = (*addr)(nil)

func NewAddr(host string) Addr {
	return &addr{host: host}
}
