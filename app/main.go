package main

import (
	"fmt"
	"github.com/codecrafters-io/dns-server-starter-go/pkg/dns"
	"net"
)

func main() {
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

		fmt.Println("DNS Message: ", buf[:size])

		m, err := dns.RawMessage(buf[:size]).Parse()
		if err != nil {
			fmt.Println("Error parsing message:", err)
			break
		}
		fmt.Printf("Start processing. ID %d\n", m.Header.ID)

		// TODO - Implement DNS server logic here.

		ipv4 := net.ParseIP("8.8.8.8").To4()
		if ipv4 == nil {
			fmt.Println("Invalid IPv4 address")
			return
		}

		fmt.Println("Parsed Message: ", m)

		rm := m.Respond(60, ipv4)

		fmt.Println("Respond Message: ", rm)
		fmt.Println("Respond Row Message: ", rm.Serialize())

		_, err = udpConn.WriteToUDP(rm.Serialize(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}

		fmt.Println("Response sent. ID", rm.Header.ID)
	}
}
