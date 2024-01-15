package dns

type Message struct {
	header    Header
	questions Questions
}

func (m *Message) Serialize() []byte {
	return append(m.header.Serialize(), m.questions.Serialize()...)
}

func NewMessage() *Message {
	return &Message{
		header: Header{
			ID: 1234,
			Flags: HeaderFlags{
				QR: 1,
			},
			QDCOUNT: 1,
			ANCOUNT: 0,
			NSCOUNT: 0,
			ARCOUNT: 0,
		},
		questions: Questions{
			Question{
				NAME:  parseDomainName("codecrafters.io"),
				TYPE:  1,
				CLASS: 1,
			},
		},
	}
}
