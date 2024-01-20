package dns

import (
	"fmt"
	"net"
	"time"
)

type Forwarder struct {
	udpAddr *net.UDPAddr
}

func NewForwarder(resolverAddress string) (*Forwarder, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", resolverAddress)
	if err != nil {
		return nil, fmt.Errorf("error resolving UDP address: %w", err)
	}
	return &Forwarder{udpAddr: udpAddr}, nil
}

func (f *Forwarder) Forward(rm RawMessage) (RawMessage, error) {
	conn, err := net.DialUDP("udp", nil, f.udpAddr)
	if err != nil {
		return RawMessage{}, fmt.Errorf("error dialing UDP: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write(rm)
	if err != nil {
		return RawMessage{}, fmt.Errorf("error sending DNS request: %w", err)
	}

	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return RawMessage{}, fmt.Errorf("error setting read deadline: %w", err)
	}

	buffer := make([]byte, 1024)
	length, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return RawMessage{}, fmt.Errorf("error reading UDP response: %w", err)
	}

	return buffer[:length], nil
}
