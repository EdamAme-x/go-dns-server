package main

import (
	"fmt"
	"net"

	"go-dns-server/dns"
)

func main() {
	fmt.Println("Starting DNS server...")

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := buf[:size]

		request, err := dns.ParseRequest(receivedData)

		if err != nil {
			fmt.Println("Failed to parse request:", err)
			continue
		}

		header := dns.CreateDNSHeader(request.Header.ID)

		header.QR = 1

		header.OPCODE = uint8(request.Header.OPCODE)

		if header.OPCODE == 0 {
			header.RCODE = dns.NoError
		} else {
			header.RCODE = dns.NotImplemented
		}

		header.RD = uint8(request.Header.RD)

		header.QDCOUNT = uint16(len(request.Questions))
		header.ANCOUNT = uint16(len(request.Questions))
		header.NSCOUNT = uint16(request.Header.NSCOUNT)
		header.ARCOUNT = uint16(request.Header.ARCOUNT)

		resp := dns.Response{}

		resp.SetHeader(header)

		for i := range request.Questions {
			question := dns.CreateQuestion(request.Questions[i].Domain)

			question.Type = dns.Record_A
			question.Class = dns.Class_IN

			resp.SetQuestion(question)
		}

		for i := range request.Questions {
			answer := dns.CreateAnswer(request.Questions[i].Domain)

			answer.Type = dns.Record_A
			answer.Class = dns.Class_IN
			answer.TTL = 60
			data := []byte{1, 1, 1, 1}
			answer.RDLENGTH = uint16(len(data))
			answer.RDATA = data

			resp.SetAnswer(answer)
		}

		response := resp.CreateResponse()

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
