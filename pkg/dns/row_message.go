package dns

import (
	"encoding/binary"
	"fmt"
)

type RowHeaderFlags []byte

type RowHeader []byte

type RowLabel []byte

type RowName []byte

type RowQuestion []byte

type RowQuestions []byte

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

func (qs RowQuestions) parse() (Questions, error) {
	var questions Questions
	offset := 0

	for offset < len(qs) {
		labelOffset, err := findLabelPointerOffset(offset, qs[offset:])
		if labelOffset.labelEndOffset > 0 {
			endOfName := labelOffset.labelEndOffset
			endOfQuestion := offset + endOfName + 5

			if endOfQuestion > len(qs) {
				return nil, fmt.Errorf("invalid question length")
			}

			questions = append(questions, RowQuestion(qs[offset:endOfQuestion]).parse())

			offset = endOfQuestion

			continue
		}
		if labelOffset.compressionPointerOffset > 0 {
			pointerOffset := labelOffset.compressionPointerOffset - 12
			endOfName, _ := findNullOffset(qs[pointerOffset:])

			var name RowName
			name = append(name, qs[offset:labelOffset.value-1]...)
			name = append(name, qs[pointerOffset:pointerOffset+endOfName+1]...)

			endOfQuestion := labelOffset.value + 5

			var rq RowQuestion
			rq = append(rq, name...)
			rq = append(rq, qs[labelOffset.value+1:endOfQuestion]...)
			questions = append(questions, rq.parse())

			offset = endOfQuestion

			continue
		}
		return Questions{}, err
	}

	return questions, nil
}

func (p RawMessage) Parse() (Message, error) {
	qs, err := RowQuestions(p[12:]).parse()
	if err != nil {
		return Message{}, err
	}
	return Message{
		Header:    RowHeader(p[0:12]).parse(),
		Questions: qs,
	}, nil
}