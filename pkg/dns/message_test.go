package dns

import (
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestHeader_Serialize(t *testing.T) {
	var expected = []byte{
		0x04, 0xD2, // ID
		0x80, 0x00, // Flags
		0x00, 0x01, // QDCOUNT
		0x00, 0x01, // ANCOUNT
		0x00, 0x00, // NSCOUNT
		0x00, 0x00, // ARCOUNT
	}

	header := Header{
		ID: 1234,
		Flags: HeaderFlags{
			QR: 1,
		},
		QDCOUNT: 1,
		ANCOUNT: 1,
		NSCOUNT: 0,
		ARCOUNT: 0,
	}

	require.Equal(t, expected, header.serialize(), "Serialized headers should match expected value")
}

func TestLabel_Serialize(t *testing.T) {
	var expected = []byte{
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	}

	var l Label = "google"

	require.Equal(t, expected, l.serialize(), "Label serialization should match expected value")
}

func TestName_Serialize(t *testing.T) {
	var expected = []byte{
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x0,
	}

	var n Name = []Label{"google", "com"}

	require.Equal(t, expected, n.serialize(), "Name serialization should match expected value")
}

func TestQuestion_Serialize(t *testing.T) {
	var expected = []byte{
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x0,
		0x00, 0x01,
		0x00, 0x01,
	}

	question := NewQuestion("google.com", 1, 1)

	require.Equal(t, expected, question.serialize(), "Serialized question should match expected value")
}

func TestQuestion_answer(t *testing.T) {
	q := NewQuestion("google.com", 1, 1)
	rdata := net.ParseIP("8.8.8.8").To4()
	a := q.answer(60, rdata)

	expected := NewAnswer(Name{"google", "com"}, 1, 1, 60, uint16(len(rdata)), rdata)

	require.Equal(t, expected, a, "Answer should match expected value")
}

func TestMessage_Serialize(t *testing.T) {
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
		Questions: Questions{
			NewQuestion("google.com", 1, 1),
		},
		Answers: Answers{
			NewAnswer(Name{"google", "com"}, 1, 1, 1, uint16(len(rdata)), rdata),
		},
	}

	require.Equal(t, expected, message.Serialize(), "Serialized headers should match expected value")
}

func TestMessage_Respond(t *testing.T) {
	rdata := net.ParseIP("8.8.8.8").To4()
	rdlen := uint16(len(rdata))
	expected := Message{
		Header: Header{
			ID: 1234,
			Flags: HeaderFlags{
				QR:    1,
				RCODE: 4,
			},
			QDCOUNT: 1,
			ANCOUNT: 1,
			NSCOUNT: 0,
			ARCOUNT: 0,
		},
		Questions: Questions{NewQuestion("google.com", 1, 1)},
		Answers:   Answers{NewAnswer(Name{"google", "com"}, 1, 1, 1, rdlen, rdata)},
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
		Questions: Questions{NewQuestion("google.com", 1, 1)},
	}

	require.Equal(t, expected, message.Respond(1, rdata), "Responded message should match expected value")
}
