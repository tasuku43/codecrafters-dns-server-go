package dns

import (
	"bytes"
	"encoding/binary"
)

type Answer struct {
	NAME    Name
	TYPE    uint16
	CLASS   uint16
	TTL     uint32
	RDLENGH uint16
	RDATA   []byte
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
