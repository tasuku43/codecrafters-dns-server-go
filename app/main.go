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

		ipv4 := net.ParseIP("8.8.8.8").To4()
		if ipv4 == nil {
			fmt.Println("Invalid IPv4 address")
			return
		}

		rm := m.Respond(60, ipv4)

		_, err = udpConn.WriteToUDP(rm.Serialize(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
