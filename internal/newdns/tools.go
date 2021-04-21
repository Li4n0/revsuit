package newdns

import (
	"math"
	"strings"
	"time"

	"github.com/miekg/dns"
)

// IsDomain returns whether the name is a valid domain and if requested also
// fully qualified.
func IsDomain(name string, fqdn bool) bool {
	_, ok := dns.IsDomainName(name)
	return ok && (!fqdn || fqdn && dns.IsFqdn(name))
}

// InZone returns whether the provided name is part of the provided zone. Will
// always return false if the provided domains are not valid.
func InZone(zone, name string) bool {
	// check domains
	if !IsDomain(zone, false) || !IsDomain(name, false) {
		return false
	}

	return dns.IsSubDomain(zone, name)
}

// TrimZone will remove the zone from the specified name.
func TrimZone(zone, name string) string {
	// return immediately if not in zone
	if !InZone(zone, name) {
		return name
	}

	// count zone labels
	count := dns.CountLabel(zone)

	// get segments
	labels := dns.SplitDomainName(name)

	// get new labels
	newLabels := labels[0 : len(labels)-count]

	// join name
	newName := strings.Join(newLabels, ".")

	return newName
}

// NormalizeDomain will normalize the provided domain name by removing space
// around the name and lowercase it if request.
func NormalizeDomain(name string, lower, makeFQDN, removeFQDN bool) string {
	// remove spaces
	name = strings.TrimSpace(name)

	// lowercase if requested
	if lower {
		name = strings.ToLower(name)
	}

	// make FQDN if requested
	if makeFQDN {
		name = dns.Fqdn(name)
	}

	// remove FQDN if requested
	if removeFQDN && dns.IsFqdn(name) {
		name = name[:len(name)-1]
	}

	return name
}

// SplitDomain will split the provided domain either in separate labels or
// hierarchical labels. The later allows walking a domain up to the root.
func SplitDomain(name string, hierarchical bool) []string {
	// normalize name
	name = NormalizeDomain(name, false, false, true)

	// return nil if empty
	if name == "" {
		return nil
	}

	// split in labels
	if !hierarchical {
		return dns.SplitDomainName(name)
	}

	// prepare list
	var list []string

	// walk domain
	for off, end := 0, false; !end; off, end = dns.NextLabel(name, off) {
		list = append(list, name[off:])
	}

	return list
}

// TransferCase will transfer the case from the source name to the destination.
// For the source "foo.AAA.com." and destination "aaa.com" the function will
// return "AAA.com". The source must be either a child or the same as the
// destination.
func TransferCase(source, destination string) string {
	// get lower variants
	lowSource := strings.ToLower(source)
	lowDestination := strings.ToLower(destination)

	// get index of destination in source
	index := strings.Index(lowSource, lowDestination)
	if index < 0 {
		return destination
	}

	// take shared part from source
	return source[index:]
}

func emailToDomain(email string) string {
	// split on at
	parts := strings.Split(email, "@")

	// replace dots in username
	parts[0] = strings.ReplaceAll(parts[0], ".", "\\.")

	// join domain
	name := parts[0] + "." + parts[1]

	return dns.Fqdn(name)
}

func toSeconds(d time.Duration) uint32 {
	return uint32(math.Ceil(d.Seconds()))
}
