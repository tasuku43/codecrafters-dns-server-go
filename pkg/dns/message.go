package dns

import "bytes"

type Message struct {
	Header   Header
	Question Question
	Answer   Answer
}

type RawMessage []byte

func (p RawMessage) Parse() Message {
	return Message{
		Header:   RowHeader(p[0:12]).parse(),
		Question: RowQuestion(p[12:]).parse(),
	}
}

func (m *Message) Serialize() []byte {
	var buffer bytes.Buffer

	buffer.Write(m.Header.Serialize())
	buffer.Write(m.Question.Serialize())
	buffer.Write(m.Answer.serialize())

	return buffer.Bytes()
}

func (m *Message) Respond(ttl uint32, rdata []byte) Message {
	rm := Message{
		Header:   m.Header,
		Question: m.Question,
		Answer:   m.Question.Answer(ttl, rdata),
	}

	rm.Header.Flags.QR = 1
	rm.Header.QDCOUNT = 1
	rm.Header.ANCOUNT = 1

	return rm
}
