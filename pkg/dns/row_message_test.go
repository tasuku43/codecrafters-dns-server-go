package dns

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

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

func TestRowLabel_Parse(t *testing.T) {
	expected := Label("google")

	data := []byte{
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	}
	actual := RowLabel(data).parse()

	require.Equal(t, expected, actual, "Parsed label should match expected value")
}

func TestRowName_Parse(t *testing.T) {
	expected := Name{"google", "com"}

	data := []byte{
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
		0x03, 0x63, 0x6f, 0x6d, 0x0,
	}
	actual := RowName(data).parse()

	require.Equal(t, expected, actual, "Parsed name should match expected value")
}

func TestRowQuestion_Parse(t *testing.T) {
	expected := Question{
		NAME:  Name{"google", "com"},
		TYPE:  1,
		CLASS: 1,
	}

	data := []byte{
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
		0x03, 0x63, 0x6f, 0x6d, 0x0,
		0x00, 0x01,
		0x00, 0x01,
	}
	actual := RowQuestion(data).parse()

	require.Equal(t, expected, actual, "Parsed question should match expected value")
}

func TestRawMessage_Parse(t *testing.T) {
	expected := Message{
		Header: Header{
			ID: 1234,
			Flags: HeaderFlags{
				QR: 1,
			},
			QDCOUNT: 2,
			ANCOUNT: 0,
			NSCOUNT: 0,
			ARCOUNT: 0,
		},
		Questions: Questions{
			Question{
				NAME:  Name{"abc", "example", "com"},
				TYPE:  1,
				CLASS: 1,
			},
			Question{
				NAME:  Name{"def", "example", "com"},
				TYPE:  1,
				CLASS: 1,
			},
		},
	}

	data := []byte{
		// Header
		0x04, 0xD2, // ID
		0x80, 0x00, // Flags
		0x00, 0x02, // QDCOUNT
		0x00, 0x00, // ANCOUNT
		0x00, 0x00, // NSCOUNT
		0x00, 0x00, // ARCOUNT
		// Question 1
		0x03, 0x61, 0x62, 0x63, // abc
		0x07, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, // example
		0x03, 0x63, 0x6f, 0x6d, // com
		0x00,       // null terminator
		0x00, 0x01, // TYPE
		0x00, 0x01, // CLASS
		// Question 2
		0x03, 0x64, 0x65, 0x66, // def
		0x07, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, // example
		0x03, 0x63, 0x6f, 0x6d, // com
		0x00,       // null terminator
		0x00, 0x01, // TYPE
		0x00, 0x01, // CLASS
	}
	actual, _ := RawMessage(data).Parse()

	require.Equal(t, expected, actual, "Parsed message should match expected value")
}

func TestRawMessage_ParseWithAnswers(t *testing.T) {
	data := []byte{
		4, 210, 129, 0, 0, 1, 0, 1, 0, 0, 0, 0, 12, 99, 111, 100, 101, 99, 114, 97, 102, 116, 101, 114, 115, 2, 105, 111, 0, 0, 1, 0, 1, 12, 99, 111, 100, 101, 99, 114, 97, 102, 116, 101, 114, 115, 2, 105, 111, 0, 0, 1, 0, 1, 0, 0, 14, 16, 0, 4, 76, 76, 21, 21,
	}

	rm := RawMessage(data)
	a, _ := rm.Parse()
	fmt.Println(a)
}
