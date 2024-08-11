package dns

type Question struct {
	Domain  string
	Type  RecordType
	Class ClassType
}

func CreateQuestion(domain string) Question {
	return Question{
		Domain:  domain,
	}
}
