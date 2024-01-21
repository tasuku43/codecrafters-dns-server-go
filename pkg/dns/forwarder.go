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

func (f *Forwarder) Forward(m Message) (Message, error) {
	conn, err := net.DialUDP("udp", nil, f.udpAddr)
	if err != nil {
		return Message{}, fmt.Errorf("error dialing UDP: %w", err)
	}
	defer conn.Close()

	messages := Messages{}
	for _, message := range m.Split() {
		_, err = conn.Write(message.Serialize())
		if err != nil {
			return Message{}, fmt.Errorf("error sending DNS request: %w", err)
		}

		err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			return Message{}, fmt.Errorf("error setting read deadline: %w", err)
		}

		buffer := make([]byte, 1024)
		length, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			return Message{}, fmt.Errorf("error reading UDP response: %w", err)
		}

		resMessage, err := RawMessage(buffer[:length]).Parse()
		if err != nil {
			return Message{}, fmt.Errorf("error parsing UDP response: %w", err)
		}

		messages = append(messages, resMessage)
	}

	return messages.Merge(), nil
}
