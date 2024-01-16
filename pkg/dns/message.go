package dns

type Message struct {
	Header   Header
	Question Question
}

type RawMessage []byte

func (p RawMessage) Parse() Message {
	return Message{
		Header:   RowHeader(p[0:12]).parse(),
		Question: RowQuestion(p[12:]).parse(),
	}
}

func (m *Message) Serialize() []byte {
	return append(m.Header.Serialize(), m.Question.Serialize()...)
}
