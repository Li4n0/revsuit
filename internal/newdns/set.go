package newdns

import (
	"time"

	"github.com/pkg/errors"
)

// Set is a set of records.
type Set struct {
	// The FQDN of the set.
	Name string

	// The type of the record.
	Type Type

	// The records in the set.
	Records []Record

	// The TTL of the set.
	//
	// Default: 5m.
	TTL time.Duration
}

// Validate will validate the set and ensure defaults.
func (s *Set) Validate() error {
	// check name
	if !IsDomain(s.Name, true) {
		return errors.Errorf("invalid name: %s", s.Name)
	}

	// check type
	if !s.Type.valid() {
		return errors.Errorf("invalid type: %d", s.Type)
	}

	// check records
	if len(s.Records) == 0 {
		return errors.Errorf("missing records")
	}

	// check CNAME records
	if s.Type == CNAME && len(s.Records) > 1 {
		return errors.Errorf("multiple CNAME records")
	}

	// validate records
	for _, record := range s.Records {
		err := record.Validate(s.Type)
		if err != nil {
			return errors.Wrap(err, "invalid record")
		}
	}

	// check for duplicate addresses if not TXT
	if len(s.Records) > 1 && s.Type != TXT {
		for i := 0; i < len(s.Records)-1; i++ {
			if s.Records[i].Address == s.Records[i+1].Address {
				return errors.Errorf("duplicate address: %s", s.Records[i].Address)
			}
		}
	}

	// set default ttl
	if s.TTL == 0 {
		s.TTL = 5 * time.Minute
	}

	return nil
}
