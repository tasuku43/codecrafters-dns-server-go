package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

type HeaderFlags struct {
	QR     uint16
	OPCODE uint16
	AA     uint16
	TC     uint16
	RD     uint16
	RA     uint16
	Z      uint16
	RCODE  uint16
}

type Header struct {
	ID      uint16
	Flags   HeaderFlags
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

type Label string

type Name []Label

type LabelOffset struct {
	value                    int
	labelEndOffset           int
	compressionPointerOffset int
}

type Question struct {
	NAME  Name
	TYPE  uint16
	CLASS uint16
}

type Questions []Question

type Answer struct {
	NAME    Name
	TYPE    uint16
	CLASS   uint16
	TTL     uint32
	RDLENGH uint16
	RDATA   []byte
}

type Answers []Answer

type Message struct {
	Header    Header
	Questions Questions
	Answers   Answers
}

func (f *HeaderFlags) toInt16() uint16 {
	var flags uint16 = 0

	flags |= f.QR << 15
	flags |= f.OPCODE << 11
	flags |= f.AA << 10
	flags |= f.TC << 9
	flags |= f.RD << 8
	flags |= f.RA << 7
	flags |= f.Z << 4
	flags |= f.RCODE

	return flags
}

func (h Header) serialize() []byte {
	res := make([]byte, 12)

	binary.BigEndian.PutUint16(res[0:], h.ID)
	binary.BigEndian.PutUint16(res[2:], h.Flags.toInt16())
	binary.BigEndian.PutUint16(res[4:], h.QDCOUNT)
	binary.BigEndian.PutUint16(res[6:], h.ANCOUNT)
	binary.BigEndian.PutUint16(res[8:], h.NSCOUNT)
	binary.BigEndian.PutUint16(res[10:], h.ARCOUNT)

	return res
}

func (l Label) serialize() []byte {
	length := uint8(len(l))

	return append([]byte{length}, []byte(l)...)
}

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

func (qs Questions) Count() uint16 {
	return uint16(len(qs))
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

func (q Question) answer(ttl uint32, rdata []byte) Answer {
	rdLength := uint16(len(rdata))
	return NewAnswer(q.NAME, q.TYPE, q.CLASS, ttl, rdLength, rdata)
}

func (a Answer) serialize() []byte {
	var buffer bytes.Buffer

	buffer.Write(a.NAME.serialize())

	binary.Write(&buffer, binary.BigEndian, a.TYPE)
	binary.Write(&buffer, binary.BigEndian, a.CLASS)
	binary.Write(&buffer, binary.BigEndian, a.TTL)
	binary.Write(&buffer, binary.BigEndian, a.RDLENGH)

	buffer.Write(a.RDATA)

	return buffer.Bytes()
}

func (as Answers) serialize() []byte {
	var serializedAnswers []byte

	for _, answer := range as {
		serializedAnswers = append(serializedAnswers, answer.serialize()...)
	}

	return serializedAnswers
}

func (as Answers) Count() uint16 {
	return uint16(len(as))
}

func (m *Message) Serialize() []byte {
	var buffer bytes.Buffer

	buffer.Write(m.Header.serialize())
	buffer.Write(m.Questions.serialize())
	buffer.Write(m.Answers.serialize())

	return buffer.Bytes()
}

func (m *Message) Respond(ttl uint32, rdata []byte) Message {
	rm := Message{
		Header:    m.Header,
		Questions: m.Questions,
		Answers:   m.Questions.answer(ttl, rdata),
	}

	rm.Header.Flags.QR = 1
	rm.Header.Flags.RCODE = 4
	rm.Header.QDCOUNT = rm.Questions.Count()
	rm.Header.ANCOUNT = rm.Answers.Count()
	rm.Header.NSCOUNT = 0
	rm.Header.ARCOUNT = 0

	return rm
}

func NewAnswer(name Name, qType uint16, qClass uint16, ttl uint32, rdlength uint16, rdata []byte) Answer {
	return Answer{
		NAME:    name,
		TYPE:    qType,
		CLASS:   qClass,
		TTL:     ttl,
		RDLENGH: rdlength,
		RDATA:   rdata,
	}
}

func findLabelPointerOffset(context int, slice []byte) (LabelOffset, error) {
	offset := 0
	for offset < len(slice) {
		b := slice[offset]
		if b == 0 {
			return LabelOffset{value: context + offset + 1, labelEndOffset: offset}, nil
		}
		if b&0xC0 == 0xC0 {
			if offset+1 >= len(slice) {
				return LabelOffset{}, fmt.Errorf("invalid pointer in label")
			}
			return LabelOffset{value: context + offset + 1, compressionPointerOffset: int(b&0x3F)<<8 + int(slice[offset+1])}, nil
		} else {
			offset++
		}
	}
	return LabelOffset{}, fmt.Errorf("null terminator not found")
}

func findNullOffset(slice []byte) (int, error) {
	for i, b := range slice {
		if b == 0 {
			return i, nil
		}
	}

	return 0, fmt.Errorf("null terminator not found")
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

func appendUint16ToSlice(slice []byte, value uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, value)
	return append(slice, bytes...)
}
