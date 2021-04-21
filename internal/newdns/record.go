package newdns

import (
	"net"

	"github.com/pkg/errors"
)

// Record holds a single DNS record.
type Record struct {
	// The target address for A, AAAA, CNAME and MX records.
	Address string

	// The priority for MX records.
	Priority int

	// The data for TXT records.
	Data []string
}

// Validate will validate the record.
func (r *Record) Validate(typ Type) error {
	// validate A address
	if typ == A {
		ip := net.ParseIP(r.Address)
		if ip == nil || ip.To4() == nil {
			return errors.Errorf("invalid IPv4 address: %s", r.Address)
		}
	}

	// validate AAAA address
	if typ == AAAA {
		ip := net.ParseIP(r.Address)
		if ip == nil || ip.To16() == nil {
			return errors.Errorf("invalid IPv6 address: %s", r.Address)
		}
	}

	// validate CNAME and MX addresses
	if typ == CNAME || typ == MX {
		if !IsDomain(r.Address, true) {
			return errors.Errorf("invalid domain name: %s", r.Address)
		}
	}

	// check TXT data
	if typ == TXT {
		if len(r.Data) == 0 {
			return errors.Errorf("missing data")
		}

		for _, data := range r.Data {
			if len(data) > 255 {
				return errors.Errorf("data too long")
			}
		}
	}

	// validate NS addresses
	if typ == NS {
		if !IsDomain(r.Address, true) {
			return errors.Errorf("invalid ns name: %s", r.Address)
		}
	}

	return nil
}
