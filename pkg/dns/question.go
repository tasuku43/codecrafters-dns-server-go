package dns

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type RowLabel []byte

type RowName []byte

type RowQuestion []byte

type RowQuestions []byte

type Label string

type Name []Label

type Question struct {
	NAME  Name
	TYPE  uint16
	CLASS uint16
}

type Questions []Question

func (qs Questions) serialize() []byte {
	var serializedQuestions []byte

	for _, question := range qs {
		serializedQuestions = append(serializedQuestions, question.serialize()...)
	}

	return serializedQuestions
}

func (qs Questions) answer(ttl uint32, rdata []byte) Answers {
	var answers Answers

	for _, question := range qs {
		answers = append(answers, question.answer(ttl, rdata))
	}

	return answers
}

func (l RowLabel) parse() Label {
	length := l[0]
	return Label(l[1 : 1+length])
}

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

func (qs RowQuestions) parse() (Questions, error) {
	var questions Questions
	offset := 0

	for offset < len(qs) {
		labelOffset, _ := findLabelPointerOffset(qs[offset:])
		if labelOffset.labelEndOffset > 0 {
			endOfName := labelOffset.labelEndOffset
			endOfQuestion := offset + endOfName + 5

			if endOfQuestion > len(qs) {
				return nil, fmt.Errorf("invalid question length")
			}

			questions = append(questions, RowQuestion(qs[offset:endOfQuestion]).parse())

			offset = endOfQuestion
		}
	}

	return questions, nil
}

type LabelOffset struct {
	labelEndOffset           int
	compressionPointerOffset int
}

func findLabelPointerOffset(slice []byte) (LabelOffset, error) {
	offset := 0
	for offset < len(slice) {
		b := slice[offset]
		if b == 0 {
			return LabelOffset{labelEndOffset: offset}, nil
		}
	}
	return LabelOffset{}, fmt.Errorf("null terminator not found")
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
