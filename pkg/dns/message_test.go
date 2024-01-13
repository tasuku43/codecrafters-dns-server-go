package dns

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSerialize(t *testing.T) {
	var expected = []byte{
		0x12, 0x34, // ID
		0x80, 0x00, // Flags
		0x00, 0x01, // QDCOUNT
		0x00, 0x01, // ANCOUNT
		0x00, 0x00, // NSCOUNT
		0x00, 0x00, // ARCOUNT
	}

	headers := Headers{
		ID: 0x1234,
		Flags: HeaderFlags{
			QR: 1,
		},
		QDCOUNT: 1,
		ANCOUNT: 1,
		NSCOUNT: 0,
		ARCOUNT: 0,
	}

	result := headers.Serialize()
	require.Equal(t, expected, result, "Serialized headers should match expected value")
}
