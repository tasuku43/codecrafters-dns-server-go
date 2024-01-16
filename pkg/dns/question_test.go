package dns

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLabelSerialize(t *testing.T) {
	var expected = []byte{
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	}

	var l Label = "google"

	require.Equal(t, expected, l.serialize(), "Label serialization should match expected value")
}

func TestNameSerialize(t *testing.T) {
	var expected = []byte{
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x0,
	}

	var n Name = []Label{"google", "com"}

	require.Equal(t, expected, n.serialize(), "Name serialization should match expected value")
}

func TestQuestionSerialize(t *testing.T) {
	var expected = []byte{
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x0,
		0x00, 0x01,
		0x00, 0x01,
	}

	question := NewQuestion("google.com", 1, 1)

	require.Equal(t, expected, question.Serialize(), "Serialized question should match expected value")
}

func TestRowLabelParse(t *testing.T) {
	expected := Label("google")

	data := []byte{
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	}
	actual := RowLabel(data).parse()

	require.Equal(t, expected, actual, "Parsed label should match expected value")
}

func TestRowNameParse(t *testing.T) {
	expected := Name{"google", "com"}

	data := []byte{
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
		0x03, 0x63, 0x6f, 0x6d, 0x0,
	}
	actual := RowName(data).parse()

	require.Equal(t, expected, actual, "Parsed name should match expected value")
}

func TestRowQuestionParse(t *testing.T) {
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
