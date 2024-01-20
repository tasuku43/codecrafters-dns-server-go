package dns

import (
	"encoding/binary"
	"fmt"
)

var headerLength = 12

type RowHeaderFlags []byte

type RowHeader []byte

type RowLabel []byte

type RowName []byte

type RowQuestion []byte

type RawMessage []byte

func (h RowHeaderFlags) parse() HeaderFlags {
	flags := binary.BigEndian.Uint16(h)

	return HeaderFlags{
		QR:     flags >> 15,
		OPCODE: (flags >> 11) & 0x0F,
		AA:     (flags >> 10) & 0x01,
		TC:     (flags >> 9) & 0x01,
		RD:     (flags >> 8) & 0x01,
		RA:     (flags >> 7) & 0x01,
		Z:      (flags >> 4) & 0x07,
		RCODE:  flags & 0x0F,
	}
}

func (h RowHeader) parse() Header {
	return Header{
		ID:      binary.BigEndian.Uint16(h[0:2]),
		Flags:   RowHeaderFlags(h[2:4]).parse(),
		QDCOUNT: binary.BigEndian.Uint16(h[4:6]),
		ANCOUNT: binary.BigEndian.Uint16(h[6:8]),
		NSCOUNT: binary.BigEndian.Uint16(h[8:10]),
		ARCOUNT: binary.BigEndian.Uint16(h[10:12]),
	}
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

func (rm RawMessage) Parse() (Message, error) {
	header := RowHeader(rm[0:headerLength]).parse()

	offset := headerLength

	// Parse questions
	var questions Questions
	for i := 0; i < int(header.QDCOUNT); i++ {
		labelOffset, err := findNameEndOrPointerOffset(offset, rm[offset:])
		if labelOffset.labelEndOffset > 0 {
			endOfName := labelOffset.labelEndOffset
			endOfQuestion := offset + endOfName + 5

			if endOfQuestion > len(rm) {
				return Message{}, fmt.Errorf("invalid question length")
			}

			questions = append(questions, RowQuestion(rm[offset:endOfQuestion]).parse())

			offset = endOfQuestion

			continue
		}
		if labelOffset.compressionPointerOffset > 0 {
			pointerOffset := labelOffset.compressionPointerOffset
			endOfName, _ := findNullOffset(rm[pointerOffset:])

			var name RowName
			name = append(name, rm[offset:labelOffset.value-1]...)
			name = append(name, rm[pointerOffset:pointerOffset+endOfName+1]...)

			endOfQuestion := labelOffset.value + 5

			var rq RowQuestion
			rq = append(rq, name...)
			rq = append(rq, rm[labelOffset.value+1:endOfQuestion]...)
			questions = append(questions, rq.parse())

			offset = endOfQuestion

			continue
		}
		return Message{}, err
	}

	// Parse answers
	var answers Answers
	for i := 0; i < int(header.ANCOUNT); i++ {
		nameEndOrPointerOffset, err := findNameEndOrPointerOffset(offset, rm[offset:])
		if nameEndOrPointerOffset.compressionPointerOffset > 0 {
			panic("not implemented")
		}
		if err != nil {
			return Message{}, err
		}
		nameEndOffset := nameEndOrPointerOffset.value
		name := RowName(rm[offset:nameEndOffset]).parse()

		typeClassOffset := nameEndOffset + 1
		qType := binary.BigEndian.Uint16(rm[typeClassOffset : typeClassOffset+2])
		qClass := binary.BigEndian.Uint16(rm[typeClassOffset+2 : typeClassOffset+4])

		ttlOffset := typeClassOffset + 4
		ttl := binary.BigEndian.Uint32(rm[ttlOffset : ttlOffset+4])

		rdLengthOffset := ttlOffset + 4
		rdLength := binary.BigEndian.Uint16(rm[rdLengthOffset : rdLengthOffset+2])

		rdata := rm[rdLengthOffset+2 : rdLengthOffset+2+int(rdLength)]

		answer := Answer{
			NAME:    name,
			TYPE:    qType,
			CLASS:   qClass,
			TTL:     ttl,
			RDLENGH: rdLength,
			RDATA:   rdata,
		}

		answers = append(answers, answer)
		offset += len(answer.serialize())
	}

	return Message{
		Header:    header,
		Questions: questions,
		Answers:   answers,
	}, nil
}

func findNameEndOrPointerOffset(context int, slice []byte) (LabelOffset, error) {
	offset := 0
	for offset < len(slice) {
		b := slice[offset]
		if b == 0 {
			return LabelOffset{value: context + offset, labelEndOffset: offset}, nil
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
