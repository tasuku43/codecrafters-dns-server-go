package dns

import "encoding/binary"

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

type Headers struct {
	ID      uint16
	Flags   HeaderFlags
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

func (h Headers) Serialize() []byte {
	res := make([]byte, 12)

	binary.BigEndian.PutUint16(res[0:], h.ID)
	binary.BigEndian.PutUint16(res[2:], h.Flags.toInt16())
	binary.BigEndian.PutUint16(res[4:], h.QDCOUNT)
	binary.BigEndian.PutUint16(res[6:], h.ANCOUNT)
	binary.BigEndian.PutUint16(res[8:], h.NSCOUNT)
	binary.BigEndian.PutUint16(res[10:], h.ARCOUNT)

	return res
}

type Message struct {
	header Headers
}

func (m *Message) Serialize() []byte {
	return m.header.Serialize()
}

func NewMessage() *Message {
	return &Message{
		header: Headers{
			ID: 1234,
			Flags: HeaderFlags{
				QR: 1,
			},
			QDCOUNT: 0,
			ANCOUNT: 0,
			NSCOUNT: 0,
			ARCOUNT: 0,
		},
	}
}
