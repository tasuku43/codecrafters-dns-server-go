package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve server address:", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Failed to connect to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter message: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if message == "exit" {
			break
		}

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Failed to send message:", err)
			continue
		}

		buf := make([]byte, 512)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Failed to read from server:", err)
			continue
		}
		fmt.Printf("Server reply: %s\n", string(buf[:n]))
	}
}
