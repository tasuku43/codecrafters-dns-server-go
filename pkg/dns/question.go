package dns

import (
	"encoding/binary"
	"strings"
)

type RowLabel []byte

type RowName []byte

type RowQuestion []byte

type Label string

func (l RowLabel) parse() Label {
	length := l[0]
	return Label(l[1 : 1+length])
}

type Name []Label

func (n RowName) parse() Name {
	var name Name
	offset := 0

	for offset < len(n) {
		length := int(n[offset])
		if length == 0 {
			break
		}
		offset++
		label := Label(n[offset : offset+length])
		name = append(name, label)
		offset += length
	}

	return name
}

type Question struct {
	NAME  Name
	TYPE  uint16
	CLASS uint16
}

func (q RowQuestion) parse() Question {
	length := len(q)
	qType := binary.BigEndian.Uint16(q[length-4 : length-2])
	qClass := binary.BigEndian.Uint16(q[length-2 : length])

	return Question{
		NAME:  RowName(q[0 : length-4]).parse(),
		TYPE:  qType,
		CLASS: qClass,
	}
}

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

func (n Name) serialize() []byte {
	var sequence []byte

	for _, label := range n {
		sequence = append(sequence, label.serialize()...)
	}

	sequence = append(sequence, 0x00)

	return sequence
}

func (q Question) serialize() []byte {
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

func (q Question) answer(ttl uint32, rdata []byte) Answer {
	rdLength := uint16(len(rdata))
	return NewAnswer(q.NAME, q.TYPE, q.CLASS, ttl, rdLength, rdata)
}
