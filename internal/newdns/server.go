package newdns

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

// Config provides configuration for a DNS server.
type Config struct {
	// The buffer size used if EDNS is enabled by a client.
	//
	// Default: 1220.
	BufferSize int

	// The list of zones handled by this server.
	//
	// Default: ["."].
	Zones []string

	// Handler is the callback that returns a zone for the specified name.
	// The returned zone must not be altered going forward.
	Handler func(name string) (*Zone, error)

	// The fallback DNS server to be used if the zones is not matched. Exact
	// zones must be provided above for this to work.
	Fallback string

	// Reporter is the callback called with request errors.
	Logger Logger
}

// Server is a DNS server.
type Server struct {
	config Config
	close  chan struct{}
}

// NewServer creates and returns a new DNS server.
func NewServer(config Config) *Server {
	// set default buffer size
	if config.BufferSize <= 0 {
		config.BufferSize = 1220
	}

	// set default zone
	if len(config.Zones) == 0 {
		config.Zones = []string{"."}
	}

	// check zones if fallback
	if config.Fallback != "" {
		for _, zone := range config.Zones {
			if zone == "." {
				panic(`fallback conflicts with the match all pattern "." (default)`)
			}
		}
	}

	return &Server{
		config: config,
		close:  make(chan struct{}),
	}
}

// Run will run a udp and tcp server on the specified address. It will return
// on the first accept error and close all servers.
func (s *Server) Run(addr string) error {
	// prepare mux
	mux := dns.NewServeMux()

	// register handler
	for _, zone := range s.config.Zones {
		mux.Handle(zone, s)
	}

	// add fallback if available
	if s.config.Fallback != "" {
		mux.Handle(".", Proxy(s.config.Fallback, s.config.Logger))
	}

	// run server
	err := Run(addr, mux, Accept(s.config.Logger), s.close)
	if err != nil {
		return err
	}

	return nil
}

// ServeDNS implements the dns.Handler interface.
func (s *Server) ServeDNS(w dns.ResponseWriter, req *dns.Msg) {
	// get question
	question := req.Question[0]

	// check class
	if question.Qclass != dns.ClassINET {
		log(s.config.Logger, Ignored, nil, nil, fmt.Sprintf("unsupported class: %s", dns.ClassToString[question.Qclass]))
		return
	}

	// log request and finish
	log(s.config.Logger, Request, req, nil, "")
	defer log(s.config.Logger, Finish, nil, nil, "")

	// prepare response
	res := new(dns.Msg)
	res.SetReply(req)

	// always compress responses
	res.Compress = true

	// set flag
	res.Authoritative = true

	// check edns
	if req.IsEdns0() != nil {
		// use edns in reply
		res.SetEdns0(uint16(s.config.BufferSize), false)

		// check version
		if req.IsEdns0().Version() != 0 {
			log(s.config.Logger, Refused, nil, nil, fmt.Sprintf("unsupported EDNS version: %d", req.IsEdns0().Version()))
			s.writeError(w, req, res, nil, dns.RcodeBadVers)
			return
		}
	}

	// check any type
	if question.Qtype == dns.TypeANY {
		log(s.config.Logger, Refused, nil, nil, "unsupported type: ANY")
		s.writeError(w, req, res, nil, dns.RcodeNotImplemented)
		return
	}

	// get name
	name := NormalizeDomain(question.Name, true, false, false)

	// get zone
	zone, err := s.config.Handler(name)
	if err != nil {
		err = errors.Wrap(err, "server handler error")
		log(s.config.Logger, BackendError, nil, err, "")
		s.writeError(w, req, res, nil, dns.RcodeServerFailure)
		return
	}

	// check zone
	if zone == nil {
		log(s.config.Logger, Refused, nil, nil, "no zone")
		res.Authoritative = false
		s.writeError(w, req, res, nil, dns.RcodeRefused)
		return
	}

	// validate zone
	err = zone.Validate()
	if err != nil {
		log(s.config.Logger, BackendError, nil, err, "")
		s.writeError(w, req, res, nil, dns.RcodeServerFailure)
		return
	}

	// answer SOA directly
	if question.Qtype == dns.TypeSOA && name == zone.Name {
		s.writeSOAResponse(w, req, res, zone)
		return
	}

	// answer NS directly
	if question.Qtype == dns.TypeNS && name == zone.Name {
		s.writeNSResponse(w, req, res, zone)
		return
	}

	// check type
	typ := Type(question.Qtype)

	// return error if type is not supported
	if !typ.valid() {
		log(s.config.Logger, Refused, nil, nil, fmt.Sprintf("unsupported type: %s", dns.TypeToString[question.Qtype]))
		s.writeError(w, req, res, zone, dns.RcodeNameError)
		return
	}

	// lookup main answer
	answer, exists, err := zone.Lookup(name, w.RemoteAddr().String(), typ)
	if err != nil {
		log(s.config.Logger, BackendError, nil, err, "")
		s.writeError(w, req, res, nil, dns.RcodeServerFailure)
		return
	}

	// check result
	if len(answer) == 0 {
		// write SOA with success code to indicate existence of other sets
		if exists {
			s.writeError(w, req, res, zone, dns.RcodeSuccess)
			return
		}

		// otherwise return name error
		s.writeError(w, req, res, zone, dns.RcodeNameError)

		return
	}

	// prepare extra set
	var extra []Set

	// TODO: Lookup glue records for NS records?

	// lookup extra sets
	for _, set := range answer {
		for _, record := range set.Records {
			switch set.Type {
			case MX:
				// lookup internal MX target A and AAAA records
				if InZone(zone.Name, record.Address) {
					ret, _, err := zone.Lookup(record.Address, w.RemoteAddr().String(), A, AAAA)
					if err != nil {
						log(s.config.Logger, BackendError, nil, err, "")
						s.writeError(w, req, res, nil, dns.RcodeServerFailure)
						return
					}

					// add to extra
					extra = append(extra, ret...)
				}
			}
		}
	}

	// set answer
	for _, set := range answer {
		res.Answer = append(res.Answer, s.convert(question.Name, zone, set)...)
	}

	// set extra
	for _, set := range extra {
		res.Extra = append(res.Extra, s.convert(question.Name, zone, set)...)
	}

	// add ns records
	for _, ns := range zone.AllNameServers {
		res.Ns = append(res.Ns, &dns.NS{
			Hdr: dns.RR_Header{
				Name:   TransferCase(question.Name, zone.Name),
				Rrtype: dns.TypeNS,
				Class:  dns.ClassINET,
				Ttl:    toSeconds(zone.NSTTL),
			},
			Ns: ns,
		})
	}

	// check if NS query
	if typ == NS {
		// move answers
		res.Ns = res.Answer
		res.Answer = nil

		// no authoritative response for other zone in NS queries
		res.Authoritative = false
	}

	// write message
	s.writeMessage(w, req, res)
}

// Close will close the server.
func (s *Server) Close() {
	defer func() { recover() }()
	close(s.close)
}

func (s *Server) writeSOAResponse(w dns.ResponseWriter, rq, rs *dns.Msg, zone *Zone) {
	// add soa record
	rs.Answer = append(rs.Answer, &dns.SOA{
		Hdr: dns.RR_Header{
			Name:   zone.Name,
			Rrtype: dns.TypeSOA,
			Class:  dns.ClassINET,
			Ttl:    toSeconds(zone.SOATTL),
		},
		Ns:      zone.MasterNameServer,
		Mbox:    emailToDomain(zone.AdminEmail),
		Serial:  1,
		Refresh: toSeconds(zone.Refresh),
		Retry:   toSeconds(zone.Retry),
		Expire:  toSeconds(zone.Expire),
		Minttl:  toSeconds(zone.MinTTL),
	})

	// add ns records
	for _, ns := range zone.AllNameServers {
		rs.Ns = append(rs.Ns, &dns.NS{
			Hdr: dns.RR_Header{
				Name:   zone.Name,
				Rrtype: dns.TypeNS,
				Class:  dns.ClassINET,
				Ttl:    toSeconds(zone.NSTTL),
			},
			Ns: ns,
		})
	}

	// write message
	s.writeMessage(w, rq, rs)
}

func (s *Server) writeNSResponse(w dns.ResponseWriter, rq, rs *dns.Msg, zone *Zone) {
	// add ns records
	for _, ns := range zone.AllNameServers {
		rs.Answer = append(rs.Answer, &dns.NS{
			Hdr: dns.RR_Header{
				Name:   zone.Name,
				Rrtype: dns.TypeNS,
				Class:  dns.ClassINET,
				Ttl:    toSeconds(zone.NSTTL),
			},
			Ns: ns,
		})
	}

	// write message
	s.writeMessage(w, rq, rs)
}

func (s *Server) writeError(w dns.ResponseWriter, rq, rs *dns.Msg, zone *Zone, code int) {
	// set code
	rs.Rcode = code

	// add soa record
	if zone != nil {
		rs.Ns = append(rs.Ns, &dns.SOA{
			Hdr: dns.RR_Header{
				Name:   zone.Name,
				Rrtype: dns.TypeSOA,
				Class:  dns.ClassINET,
				Ttl:    toSeconds(zone.SOATTL),
			},
			Ns:      zone.MasterNameServer,
			Mbox:    emailToDomain(zone.AdminEmail),
			Serial:  1,
			Refresh: toSeconds(zone.Refresh),
			Retry:   toSeconds(zone.Retry),
			Expire:  toSeconds(zone.Expire),
			Minttl:  toSeconds(zone.MinTTL),
		})
	}

	// write message
	s.writeMessage(w, rq, rs)
}

func (s *Server) writeMessage(w dns.ResponseWriter, rq, rs *dns.Msg) {
	// get buffer size
	var buffer = 512
	if rq.IsEdns0() != nil {
		buffer = int(rq.IsEdns0().UDPSize())
	}

	// determine if client is using UDP
	isUDP := w.RemoteAddr().Network() == "udp"

	// truncate message if client is using UDP and message is too long
	if isUDP && rs.Len() > buffer {
		rs.Truncated = true
		rs.Answer = nil
		rs.Ns = nil
		rs.Extra = nil
	}

	// write message
	err := w.WriteMsg(rs)
	if err != nil {
		log(s.config.Logger, NetworkError, nil, err, "")
		_ = w.Close()
		return
	}

	// log response
	log(s.config.Logger, Response, rs, nil, "")
}

func (s *Server) convert(query string, zone *Zone, set Set) []dns.RR {
	// prepare header
	header := dns.RR_Header{
		Name:   TransferCase(query, set.Name),
		Rrtype: uint16(set.Type),
		Class:  dns.ClassINET,
		Ttl:    toSeconds(set.TTL),
	}

	// ensure zone min TTL
	if set.TTL < zone.MinTTL {
		header.Ttl = toSeconds(zone.MinTTL)
	}

	// prepare list
	var list []dns.RR

	// add records
	for _, record := range set.Records {
		// construct record
		switch set.Type {
		case A:
			list = append(list, &dns.A{
				Hdr: header,
				A:   net.ParseIP(record.Address),
			})
		case AAAA:
			list = append(list, &dns.AAAA{
				Hdr:  header,
				AAAA: net.ParseIP(record.Address),
			})
		case CNAME:
			list = append(list, &dns.CNAME{
				Hdr:    header,
				Target: dns.Fqdn(record.Address),
			})
		case MX:
			list = append(list, &dns.MX{
				Hdr:        header,
				Preference: uint16(record.Priority),
				Mx:         dns.Fqdn(record.Address),
			})
		case TXT:
			list = append(list, &dns.TXT{
				Hdr: header,
				Txt: record.Data,
			})
		case NS:
			list = append(list, &dns.NS{
				Hdr: header,
				Ns:  dns.Fqdn(record.Address),
			})
		}
	}

	return list
}
