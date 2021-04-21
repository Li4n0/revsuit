package newdns

import "github.com/miekg/dns"

// Type denotes the DNS record type.
type Type uint16

const (
	// A records return IPV4 addresses.
	A = Type(dns.TypeA)

	// AAAA records return IPV6 addresses.
	AAAA = Type(dns.TypeAAAA)

	// CNAME records return other DNS names.
	CNAME = Type(dns.TypeCNAME)

	// MX records return mails servers with their priorities. The target mail
	// servers must itself be returned with an A or AAAA record.
	MX = Type(dns.TypeMX)

	// TXT records return arbitrary text data.
	TXT = Type(dns.TypeTXT)

	// NS records delegate names to other name servers.
	NS = Type(dns.TypeNS)

	// REBINDING records return different ip when the same client requests twice.
	REBINDING = Type(99)
)

func (t Type) valid() bool {
	switch t {
	case A, AAAA, CNAME, MX, TXT, NS, REBINDING:
		return true
	default:
		return false
	}
}

func typeInList(list []Type, needle Type) bool {
	for _, t := range list {
		if t == needle {
			return true
		}
	}

	return false
}
