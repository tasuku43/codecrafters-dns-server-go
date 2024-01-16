package dns

import "encoding/binary"

type RowHeader []byte

type RowHeaderFlags []byte

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

type Header struct {
	ID      uint16
	Flags   HeaderFlags
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

func (h Header) Serialize() []byte {
	res := make([]byte, 12)

	binary.BigEndian.PutUint16(res[0:], h.ID)
	binary.BigEndian.PutUint16(res[2:], h.Flags.toInt16())
	binary.BigEndian.PutUint16(res[4:], h.QDCOUNT)
	binary.BigEndian.PutUint16(res[6:], h.ANCOUNT)
	binary.BigEndian.PutUint16(res[8:], h.NSCOUNT)
	binary.BigEndian.PutUint16(res[10:], h.ARCOUNT)

	return res
}
