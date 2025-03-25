package networking

type DualClient struct {
	quitCh chan struct{}
	TCP    *TCPClient
	UDP    *UDPClient
}

func (c *DualClient) Connect() error {
	err := c.TCP.Connect()
	if err != nil {
		return err
	}
	go c.TCP.Listen()

	err = c.UDP.Connect()
	if err != nil {
		return err
	}
	go c.UDP.Listen()

	return nil
}

func (c *DualClient) Close() {
	c.TCP.Close()
	c.UDP.Close()
}

func NewClient(serverAddr string) *DualClient {
	return &DualClient{
		quitCh: make(chan struct{}),
		TCP:    NewTCPClient(serverAddr),
		UDP:    NewUDPClient(serverAddr),
	}
}
