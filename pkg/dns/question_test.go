package dns

import (
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

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

func TestRowQuestions_Parse(t *testing.T) {
	expected := Questions{
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
	}

	data := []byte{
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

	actual, _ := RowQuestions(data).parse()

	require.Equal(t, expected, actual, "Parsed questions should match expected value")
}

func TestQuestion_answer(t *testing.T) {
	q := NewQuestion("google.com", 1, 1)
	rdata := net.ParseIP("8.8.8.8").To4()
	a := q.answer(60, rdata)

	expected := NewAnswer(Name{"google", "com"}, 1, 1, 60, uint16(len(rdata)), rdata)

	require.Equal(t, expected, a, "Answer should match expected value")
}
