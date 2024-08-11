package dns

import "encoding/binary"

type Response struct {
	Header   []byte
	Question []byte
	Answer   []byte
}

func (r *Response) SetHeader(header Header) {
	r.Header = make([]byte, 12)

	binary.BigEndian.PutUint16(r.Header[0:2], header.ID)

	binary.BigEndian.PutUint16(r.Header[2:4], combineFlags(uint(header.QR), uint(header.OPCODE), uint(header.AA), uint(header.TC), uint(header.RD), uint(header.RA), uint(header.Z), uint(header.RCODE)))
	binary.BigEndian.PutUint16(r.Header[4:6], header.QDCOUNT)
	binary.BigEndian.PutUint16(r.Header[6:8], header.ANCOUNT)
	binary.BigEndian.PutUint16(r.Header[8:10], header.NSCOUNT)
	binary.BigEndian.PutUint16(r.Header[10:12], header.ARCOUNT)
}


func (r *Response) SetQuestion(question Question) {
	r.Question = []byte{}

	r.Question = append(r.Question, encodeDomain(question.Domain)...)
	r.Question = AppendUint16(r.Question, (uint16(question.Type)))
	r.Question = AppendUint16(r.Question, (uint16(question.Class)))
}

func (r *Response) SetAnswer(answer Answer) {
	if r.Answer == nil {
		r.Answer = []byte{}
	}
	
	r.Answer = append(r.Answer, encodeDomain(answer.Domain)...)
	r.Answer = AppendUint16(r.Answer, (uint16(answer.Type)))
	r.Answer = AppendUint16(r.Answer, (uint16(answer.Class)))
	r.Answer = AppendUint32(r.Answer, answer.TTL)
	r.Answer = AppendUint16(r.Answer, answer.RDLENGTH)
	r.Answer = append(r.Answer, answer.RDATA...)
}

func (r *Response) CreateResponse() []byte {
	return append(r.Header, append(r.Question, r.Answer...)...)
}

func AppendUint16(b []byte, v uint16) []byte {
	return append(b,
		byte(v>>8),
		byte(v),
	)
}

func AppendUint32(b []byte, v uint32) []byte {
	return append(b,
		byte(v>>24),
		byte(v>>16),
		byte(v>>8),
		byte(v),
	)
}

func combineFlags(qr, opcode, aa, tc, rd, ra, z, rcode uint) uint16 {
	return uint16(qr<<15 | opcode<<11 | aa<<10 | tc<<9 | rd<<8 | ra<<7 | z<<4 | rcode)
}
