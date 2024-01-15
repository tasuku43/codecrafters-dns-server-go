package dns

import (
	"encoding/binary"
	"strings"
)

type Label string

type Name []Label

type Question struct {
	NAME  Name
	TYPE  uint16
	CLASS uint16
}

type Questions []Question

func (l Label) serialize() []byte {
	length := uint8(len(l))

	return append([]byte{length}, []byte(l)...)
}

func parseDomainName(n string) Name {
	labels := strings.Split(n, ".")
	name := make(Name, len(labels))

	for i, label := range labels {
		name[i] = Label(label)
	}

	return name
}

func NewQuestion(n string, t uint16, c uint16) Question {
	return Question{
		NAME:  parseDomainName(n),
		TYPE:  t,
		CLASS: c,
	}
}

func (q Questions) Serialize() []byte {
	var serializedQuestions []byte

	for _, question := range q {
		serializedQuestions = append(serializedQuestions, question.Serialize()...)
	}

	return serializedQuestions
}

func (n Name) serialize() []byte {
	var sequence []byte

	for _, label := range n {
		sequence = append(sequence, label.serialize()...)
	}

	sequence = append(sequence, 0x00)

	return sequence
}

func (q Question) Serialize() []byte {
	var serializedQuestion []byte

	serializedQuestion = append(serializedQuestion, q.NAME.serialize()...)
	serializedQuestion = appendUint16ToSlice(serializedQuestion, q.TYPE)
	serializedQuestion = appendUint16ToSlice(serializedQuestion, q.CLASS)

	return serializedQuestion
}

func appendUint16ToSlice(slice []byte, value uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, value)
	return append(slice, bytes...)
}
