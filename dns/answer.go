package dns

type Answer struct {
	Domain  string
	Type  RecordType
	Class ClassType
	TTL   uint32
	RDLENGTH uint16
	RDATA []byte
}

func CreateAnswer(domain string) Answer {
	return Answer{
		Domain:  domain,
	}
}
