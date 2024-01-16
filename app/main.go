package main

import (
	"fmt"
	"github.com/codecrafters-io/dns-server-starter-go/pkg/dns"
	"net"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		m := dns.RawMessage(buf[:size]).Parse()

		// TODO - Implement DNS server logic here.

		m.Header.Flags.QR = 1
		m.Header.QDCOUNT = 1
		m.Question.TYPE = 1
		m.Question.CLASS = 1

		_, err = udpConn.WriteToUDP(m.Serialize(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
