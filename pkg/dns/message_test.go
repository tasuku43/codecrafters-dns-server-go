package dns

import (
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestMessageSerialize(t *testing.T) {
	var expected = []byte{
		// Header
		0x04, 0xD2, // ID
		0x80, 0x00, // Flags
		0x00, 0x01, // QDCOUNT
		0x00, 0x01, // ANCOUNT
		0x00, 0x00, // NSCOUNT
		0x00, 0x00, // ARCOUNT
		// Question
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x0,
		0x00, 0x01,
		0x00, 0x01,
		// Answer
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x0,
		0x00, 0x01,
		0x00, 0x01,
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x04,
		0x08, 0x08, 0x08, 0x08,
	}

	rdata := net.ParseIP("8.8.8.8").To4()
	message := Message{
		Header: Header{
			ID: 1234,
			Flags: HeaderFlags{
				QR: 1,
			},
			QDCOUNT: 1,
			ANCOUNT: 1,
			NSCOUNT: 0,
			ARCOUNT: 0,
		},
		Question: NewQuestion("google.com", 1, 1),
		Answer:   NewAnswer(Name{"google", "com"}, 1, 1, 1, uint16(len(rdata)), rdata),
	}

	require.Equal(t, expected, message.Serialize(), "Serialized headers should match expected value")
}

func TestRawMessage_Parse(t *testing.T) {
	expected := Message{
		Header: Header{
			ID: 1234,
			Flags: HeaderFlags{
				QR: 1,
			},
			QDCOUNT: 1,
			ANCOUNT: 1,
			NSCOUNT: 0,
			ARCOUNT: 0,
		},
		Question: NewQuestion("google.com", 1, 1),
	}

	data := []byte{
		// Header
		0x04, 0xD2, // ID
		0x80, 0x00, // Flags
		0x00, 0x01, // QDCOUNT
		0x00, 0x01, // ANCOUNT
		0x00, 0x00, // NSCOUNT
		0x00, 0x00, // ARCOUNT
		// Question
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x0,
		0x00, 0x01,
		0x00, 0x01,
	}
	actual := RawMessage(data).Parse()

	require.Equal(t, expected, actual, "Parsed message should match expected value")
}

func TestMessageRespond(t *testing.T) {
	rdata := net.ParseIP("8.8.8.8").To4()
	rdlen := uint16(len(rdata))
	expected := Message{
		Header: Header{
			ID: 1234,
			Flags: HeaderFlags{
				QR: 1,
			},
			QDCOUNT: 1,
			ANCOUNT: 1,
			NSCOUNT: 0,
			ARCOUNT: 0,
		},
		Question: NewQuestion("google.com", 1, 1),
		Answer:   NewAnswer(Name{"google", "com"}, 1, 1, 1, rdlen, rdata),
	}

	message := Message{
		Header: Header{
			ID: 1234,
			Flags: HeaderFlags{
				QR: 1,
			},
			QDCOUNT: 0,
			ANCOUNT: 0,
			NSCOUNT: 0,
			ARCOUNT: 0,
		},
		Question: NewQuestion("google.com", 1, 1),
	}

	require.Equal(t, expected, message.Respond(1, rdata), "Responded message should match expected value")
}
