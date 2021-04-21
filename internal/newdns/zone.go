package newdns

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Zone describes a single authoritative DNS zone.
type Zone struct {
	// The FQDN of the zone e.g. "example.com.".
	Name string

	// The FQDN of the master mame server responsible for this zone. The FQDN
	// must be returned as A and AAAA record by the parent zone.
	MasterNameServer string

	// A list of FQDNs to all authoritative name servers for this zone. The
	// FQDNs must be returned as A and AAAA records by the parent zone. It is
	// required to announce at least two distinct name servers per zone.
	AllNameServers []string

	// The email address of the administrator e.g. "hostmaster@example.com".
	//
	// Default: "hostmaster@NAME".
	AdminEmail string

	// The refresh interval.
	//
	// Default: 6h.
	Refresh time.Duration

	// The retry interval for the zone.
	//
	// Default: 1h.
	Retry time.Duration

	// The expiration interval of the zone.
	//
	// Default: 72h.
	Expire time.Duration

	// The TTL for the SOA record.
	//
	// Default: 15m.
	SOATTL time.Duration

	// The TTL for NS records.
	//
	// Default: 48h.
	NSTTL time.Duration

	// The minimum TTL for all records. Either this value, or the SOATTL if lower,
	// is used to determine the "negative caching TTL" which is the duration
	// caches are allowed to cache missing records (NXDOMAIN).
	//
	// Default: 5min.
	MinTTL time.Duration

	// The handler that responds to requests for this zone. The returned sets
	// must not be altered going forward.
	Handler func(name, remoteAddr string) ([]Set, error)
}

// Validate will validate the zone and ensure the documented defaults.
func (z *Zone) Validate() error {
	// check name
	if !IsDomain(z.Name, true) {
		return errors.Errorf("name not fully qualified: %s", z.Name)
	}

	// check master name server
	if !IsDomain(z.MasterNameServer, true) {
		return errors.Errorf("master server not full qualified: %s", z.MasterNameServer)
	}

	// check name server count
	if len(z.AllNameServers) < 1 {
		return errors.Errorf("missing name servers")
	}

	// check name servers
	var includesMaster bool
	for _, ns := range z.AllNameServers {
		if !IsDomain(ns, true) {
			return errors.Errorf("name server not fully qualified: %s", ns)
		}

		if ns == z.MasterNameServer {
			includesMaster = true
		}
	}

	// check master inclusion
	if !includesMaster {
		return errors.Errorf("master name server not listed as name server: %s", z.MasterNameServer)
	}

	// set default admin email
	if z.AdminEmail == "" {
		z.AdminEmail = fmt.Sprintf("hostmaster@%s", z.Name)
	}

	// check admin email
	if !IsDomain(emailToDomain(z.AdminEmail), true) {
		return errors.Errorf("admin email cannot be converted to a domain name: %s", z.AdminEmail)
	}

	// set default refresh
	if z.Refresh == 0 {
		z.Refresh = 6 * time.Hour
	}

	// set default retry
	if z.Retry == 0 {
		z.Retry = time.Hour
	}

	// set default expire
	if z.Expire == 0 {
		z.Expire = 72 * time.Hour
	}

	// set default SOA TTL
	if z.SOATTL == 0 {
		z.SOATTL = 15 * time.Minute
	}

	// set default NS TTL
	if z.NSTTL == 0 {
		z.NSTTL = 48 * time.Hour
	}

	// set default min TTL
	if z.MinTTL == 0 {
		z.MinTTL = 5 * time.Minute
	}

	// check retry
	if z.Retry >= z.Refresh {
		return errors.Errorf("retry must be less than refresh: %d", z.Retry)
	}

	// check expire
	if z.Expire < z.Refresh+z.Retry {
		return errors.Errorf("expire must be bigger than the sum of refresh and retry: %d", z.Expire)
	}

	return nil
}

// Lookup will lookup the specified name in the zone and return results for the
// specified record types. If no results are returned, the second return value
// indicates if there are other results for the specified name.
func (z *Zone) Lookup(name, remoteAddr string, needle ...Type) ([]Set, bool, error) {
	// check name
	if !IsDomain(name, true) {
		return nil, false, errors.Errorf("invalid name: %s", name)
	}

	// normalize name
	name = NormalizeDomain(name, true, false, false)

	// check name
	if !InZone(z.Name, name) {
		return nil, false, errors.Errorf("name does not belong to zone: %s", name)
	}

	// prepare result
	var result []Set

	for i := 0; ; i++ {
		// get sets
		sets, err := z.Handler(TrimZone(z.Name, name),remoteAddr)
		if err != nil {
			return nil, false, errors.Wrap(err, "zone handler error")
		}

		// return immediately if initial set is empty
		if i == 0 && len(sets) == 0 {
			return nil, false, nil
		}

		// prepare counters
		counters := map[Type]int{
			A:     0,
			AAAA:  0,
			CNAME: 0,
			MX:    0,
			TXT:   0,
		}

		// validate sets
		for _, set := range sets {
			// validate set
			err = set.Validate()
			if err != nil {
				return nil, false, errors.Wrap(err, "invalid set")
			}

			// check relationship
			if !InZone(z.Name, set.Name) {
				return nil, false, errors.Errorf("set does not belong to zone: %s", set.Name)
			}

			// increment counter
			counters[set.Type]++
		}

		// check counters
		for _, counter := range counters {
			if counter > 1 {
				return nil, false, errors.New("multiple sets for same type")
			}
		}

		// check apex CNAME
		if counters[CNAME] > 0 && name == z.Name {
			return nil, false, errors.Errorf("invalid CNAME set at apex: %s", name)
		}

		// check CNAME is stand-alone
		if counters[CNAME] > 0 && (len(sets) > 1) {
			return nil, false, errors.Errorf("other sets with CNAME set: %s", name)
		}

		// check if CNAME and query is not CNAME
		if counters[CNAME] > 0 && !typeInList(needle, CNAME) {
			// add set to result
			result = append(result, sets[0])

			// get normalized address
			address := NormalizeDomain(sets[0].Records[0].Address, true, false, false)

			// continue lookup with CNAME address if address is in zone
			if InZone(z.Name, address) {
				name = address
				continue
			}

			return result, false, nil
		}

		// add matching set
		for _, set := range sets {
			if typeInList(needle, set.Type) {
				// add set to result
				result = append(result, set)

				break
			}
		}

		// return if there are not matches, but indicate that there are sets
		// available for other types
		if len(result) == 0 {
			return nil, true, nil
		}

		return result, false, nil
	}
}
