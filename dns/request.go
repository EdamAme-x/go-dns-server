package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

type DNSMessage struct {
	Header    Header
	Questions []Question
	Answers   []Answer
}

func parseHeader(data []byte) (*Header, error) {
	if len(data) < 12 { 
		return nil, fmt.Errorf("header too short")
	}
	header := &Header{
		ID: binary.BigEndian.Uint16(data[:2]),
	}
	flags := binary.BigEndian.Uint16(data[2:4])
	header.QR = uint8(flags >> 15 & 0x01)
	header.OPCODE = uint8(flags >> 11 & 0x0F)
	header.AA = uint8(flags >> 10 & 0x01)
	header.TC = uint8(flags >> 9 & 0x01)
	header.RD = uint8(flags >> 8 & 0x01)
	header.RA = uint8(flags >> 7 & 0x01)
	header.Z = uint8(flags >> 4 & 0x07) 
	header.RCODE = ResultCode(flags & 0x0F)
	header.QDCOUNT = binary.BigEndian.Uint16(data[4:6])
	header.ANCOUNT = binary.BigEndian.Uint16(data[6:8])
	header.NSCOUNT = binary.BigEndian.Uint16(data[8:10])
	header.ARCOUNT = binary.BigEndian.Uint16(data[10:12])
	return header, nil
}

func parseDomain(buf []byte, source []byte) string {
	offset := 0
	labels := []string{}
	for {
		if buf[offset] == 0 {
			break
		}
		if (buf[offset]&0xC0)>>6 == 0b11 {
			ptr := int(binary.BigEndian.Uint16(buf[offset:offset+2]) << 2 >> 2)
			length := bytes.Index(source[ptr:], []byte{0})
			labels = append(labels, parseDomain(source[ptr:ptr+length+1], source))
			offset += 2
			continue
		}
		length := int(buf[offset])
		substring := buf[offset+1 : offset+1+length]
		labels = append(labels, string(substring))
		offset += length + 1
	}
	return strings.Join(labels, ".")
}

func ParseRequest(data []byte) (DNSMessage, error) {
	var msg DNSMessage
	var err error
	reader := bytes.NewReader(data)
	// Parse Header
	hdr, err := parseHeader(data)
	if err != nil {
		return msg, err
	}
	msg.Header = *hdr

	_, _ = reader.Seek(12, io.SeekStart)
	for i := 0; i < int(msg.Header.QDCOUNT); i++ {
		var q Question
		err = binary.Read(reader, binary.BigEndian, &q.Type)
		if err != nil {
			return msg, fmt.Errorf("failed to parse question type: %v", err)
		}
		err = binary.Read(reader, binary.BigEndian, &q.Class)
		if err != nil {
			return msg, fmt.Errorf("failed to parse question class: %v", err)
		}
		msg.Questions = append(msg.Questions, q)
	}

	for i := 0; i < int(msg.Header.ANCOUNT); i++ {
		var a Answer
		err = binary.Read(reader, binary.BigEndian, &a.Type)
		if err != nil {
			return msg, fmt.Errorf("failed to parse answer type: %v", err)
		}
		err = binary.Read(reader, binary.BigEndian, &a.Class)
		if err != nil {
			return msg, fmt.Errorf("failed to parse answer class: %v", err)
		}
		err = binary.Read(reader, binary.BigEndian, &a.TTL)
		if err != nil {
			return msg, fmt.Errorf("failed to parse answer TTL: %v", err)
		}
		err = binary.Read(reader, binary.BigEndian, &a.RDLENGTH)
		if err != nil {
			return msg, fmt.Errorf("failed to parse answer RDLENGTH: %v", err)
		}
		a.RDATA = make([]byte, a.RDLENGTH)
		_, err = reader.Read(a.RDATA)
		if err != nil {
			return msg, fmt.Errorf("failed to parse answer RDATA: %v", err)
		}
		msg.Answers = append(msg.Answers, a)
	}

	domainList := getDomains(data, msg.Header.QDCOUNT)

	for i := 0; i < len(domainList); i++ {
		msg.Questions[i].Domain = domainList[i]
	}

	return msg, nil
}

func getDomains(serializedBuf []byte, numQues uint16) []string {
	offset := 12
	domainList := []string{}
	for i := uint16(0); i < numQues; i++ {
		len := bytes.Index(serializedBuf[offset:], []byte{0})
		domain := parseDomain(serializedBuf[offset:offset+len+1], serializedBuf)
		domainList = append(domainList, domain)
		offset += len + 1
		offset += 4
	}
	return domainList
}
