package main

import (
	"fmt"

	"github.com/ouijan/ingenuity/pkg/networking"
)

func main() {
	tcpServer := networking.NewTCPServer("localhost:8080", 10)

	go func() {
		for msg := range tcpServer.MsgCh {
			fmt.Printf("From %s: %s\n", msg.Conn.RemoteAddr(), msg.Payload)
			msg.Conn.Write([]byte("Acknowledged"))
		}
	}()

	tcpServer.Start()
}
