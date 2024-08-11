package dns

import (
	"strings"
)


func encodeDomain(domain string) []byte {
	splittedDomain := strings.Split(domain, ".")
	encodedDomain := []byte{}

	for _, part := range splittedDomain {
		encodedDomain = append(encodedDomain, byte(len(part)))
		encodedDomain = append(encodedDomain, []byte(part)...)
	}

	encodedDomain = append(encodedDomain, 0x00)

	return encodedDomain
}