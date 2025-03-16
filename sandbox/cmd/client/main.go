package main

import (
	"fmt"
	"time"

	"github.com/ouijan/ingenuity/pkg/networking"
)

func main() {
	tcpClient := networking.NewTCPClient("localhost:8080")

	go func() {
		for msg := range tcpClient.MsgCh {
			fmt.Printf("From %s: %s\n", msg.Conn.RemoteAddr(), msg.Payload)
		}
	}()

	tcpClient.Connect()
	defer tcpClient.Close()

	go tcpClient.ReadLoop()

	err := tcpClient.Write([]byte("Hello, World!"))
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(1 * time.Second)
}
