package dns

import (
	"github.com/stretchr/testify/require"
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

	require.Equal(t, expected, header.Serialize(), "Serialized headers should match expected value")
}

func TestRowHeader_parse(t *testing.T) {
	expected := Header{
		ID: 1234,
		Flags: HeaderFlags{
			QR: 1,
		},
		QDCOUNT: 1,
		ANCOUNT: 1,
		NSCOUNT: 0,
		ARCOUNT: 0,
	}

	data := []byte{
		// Header
		0x04, 0xD2, // ID
		0x80, 0x00, // Flags
		0x00, 0x01, // QDCOUNT
		0x00, 0x01, // ANCOUNT
		0x00, 0x00, // NSCOUNT
		0x00, 0x00, // ARCOUNT
	}
	actual := RowHeader(data).parse()

	require.Equal(t, expected, actual, "Parsed header should match expected value")
}
