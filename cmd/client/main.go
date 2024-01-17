package main

import (
	"encoding/binary"
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

	dnsRequest := createDNSRequest([]string{"abc.example.com", "def.example.com"})

	_, err = conn.Write(dnsRequest)
	if err != nil {
		fmt.Println("Failed to send message:", err)
		return
	}

	buf := make([]byte, 512)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("Failed to read from server:", err)
		return
	}
	fmt.Printf("Server reply: %s\n", string(buf[:n]))
}

func createDNSRequest(domains []string) []byte {
	var header [12]byte
	binary.BigEndian.PutUint16(header[0:2], 1234)                 // ID
	binary.BigEndian.PutUint16(header[4:6], uint16(len(domains))) // QDCOUNT

	question := make([]byte, 0)
	for _, domain := range domains {
		question = append(question, createQuestion(domain, 1, 1)...)
	}

	return append(header[:], question...)
}

func createQuestion(domain string, qtype, qclass uint16) []byte {
	question := make([]byte, 0)
	for _, part := range strings.Split(domain, ".") {
		question = append(question, byte(len(part)))
		question = append(question, part...)
	}
	question = append(question, 0)
	question = appendUint16ToSlice(question, qtype)
	question = appendUint16ToSlice(question, qclass)
	return question
}

func appendUint16ToSlice(slice []byte, value uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, value)
	return append(slice, bytes...)
}
