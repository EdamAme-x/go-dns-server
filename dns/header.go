package dns

type Header struct {
	ID      uint16
	QR      uint8
	OPCODE 	uint8
	AA      uint8
	TC      uint8
	RD      uint8
	RA      uint8
	Z       uint8
	RCODE   ResultCode
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

type ResultCode uint8

const (
	NoError ResultCode = iota
	FormatError
	ServerFailure
	NameError
	NotImplemented
	Refused
)

func (r ResultCode) Get(resultCode ResultCode) int8 {
	switch resultCode {
	case NoError:
		return 0
	case FormatError:
		return 1
	case ServerFailure:
		return 2
	case NameError:
		return 3
	case NotImplemented:
		return 4
	case Refused:
		return 5
	}
	return 0
}

func CreateDNSHeader(id uint16) Header {
	return Header{
		ID:      id,
		QR:      0,
		OPCODE:  0,
		AA:      0,
		TC:      0,
		RD:      0,
		RA:      0,
		Z:       0,
		RCODE:   NoError,
		QDCOUNT: 0,
		ANCOUNT: 0,
		NSCOUNT: 0,
		ARCOUNT: 0,
	}
}
