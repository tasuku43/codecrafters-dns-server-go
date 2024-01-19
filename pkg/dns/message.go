package dns

import "bytes"

type Message struct {
	Header    Header
	Questions Questions
	Answers   Answers
}

type RawMessage []byte

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
