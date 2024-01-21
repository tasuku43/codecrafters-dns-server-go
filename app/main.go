package main

import (
	"flag"
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

	resolverAddress := flag.String("resolver", "", "DNS resolver address.")
	flag.Parse()

	if *resolverAddress == "" {
		fmt.Println("Error: Resolver address is required.")
		return
	}

	forwarder, err := dns.NewForwarder(*resolverAddress)
	if err != nil {
		fmt.Println("Failed to create forwarder:", err)
		return
	}

	fmt.Printf("Using DNS resolver at address: %s\n", *resolverAddress)

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		m, err := dns.RawMessage(buf[:size]).Parse()
		if err != nil {
			fmt.Println("Error parsing message:", err)
			break
		}
		fmt.Printf("Start processing. ID %d\n", m.Header.ID)
		fmt.Println("[", m.Header.ID, "]DNS Message: ", buf[:size])

		rm, err := forwarder.Forward(m)
		if err != nil {
			fmt.Println("Error forwarding message:", err)
			break
		}

		fmt.Printf("[ %d ]Response Message Header: %+v\n", m.Header.ID, rm.Header)
		fmt.Printf("[ %d ]Response Message Questions: %+v\n", m.Header.ID, rm.Questions)
		fmt.Printf("[ %d ]Response Message Answers: %+v\n", m.Header.ID, rm.Answers)

		_, err = udpConn.WriteToUDP(rm.Serialize(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}

		fmt.Println("Response sent. ID", m.Header.ID)
	}
}
