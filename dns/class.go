package dns

type ClassType uint16

const (
	_ ClassType = iota
	Class_IN
	Class_CS
	Class_CH
	Class_HS
)

/**
  IN - 1 the Internet
  CS - 2 the CSNET class (Obsolete - used only for examples in some obsolete RFCs)
  CH - 3 the CHAOS class
  HS - 4 Hesiod [Dyer 87]
*/
