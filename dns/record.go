package dns

type RecordType uint16

const (
	_ RecordType = iota
	Record_A
	Record_NS
	Record_MD
	Record_MF
	Record_CNAME
	Record_SOA
	Record_MB
	Record_MG
	Record_MR
	Record_NULL
	Record_WKS
	Record_PTR
	Record_HINFO
	Record_MINFO
	Record_MX
	Record_TXT
)

/**
	A - 1 a host address
    NS - 2 an authoritative name server
    MD - 3 a mail destination (Obsolete - use MX)
    MF - 4 a mail forwarder (Obsolete - use MX)
    CNAME - 5 the canonical name for an alias
    SOA - 6 marks the start of a zone of authority
    MB - 7 a mailbox domain name (EXPERIMENTAL)
    MG - 8 a mail group member (EXPERIMENTAL)
    MR - 9 a mail rename domain name (EXPERIMENTAL)
    NULL - 10 a null RR (EXPERIMENTAL)
    WKS - 11 a well known service description
    PTR - 12 a domain name pointer
    HINFO - 13 host information
    MINFO - 14 mailbox or mail list information
    MX - 15 mail exchange
    TXT - 16 text strings
*/
