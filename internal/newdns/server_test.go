package newdns

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

var awsNS = []string{
	"ns-1071.awsdns-05.org.",
	"ns-140.awsdns-17.com.",
	"ns-1978.awsdns-55.co.uk.",
	"ns-812.awsdns-37.net.",
}

const awsPrimaryNS = "ns-140.awsdns-17.com."

var awsOtherNS = []string{
	"ns-1074.awsdns-06.org.",
	"ns-1631.awsdns-11.co.uk.",
	"ns-243.awsdns-30.com.",
	"ns-869.awsdns-44.net.",
}

var nsRRs = []dns.RR{
	&dns.NS{
		Hdr: dns.RR_Header{
			Name:     "newdns.256dpi.com.",
			Rrtype:   dns.TypeNS,
			Class:    dns.ClassINET,
			Ttl:      172800,
			Rdlength: 23,
		},
		Ns: awsNS[0],
	},
	&dns.NS{
		Hdr: dns.RR_Header{
			Name:     "newdns.256dpi.com.",
			Rrtype:   dns.TypeNS,
			Class:    dns.ClassINET,
			Ttl:      172800,
			Rdlength: 19,
		},
		Ns: awsNS[1],
	},
	&dns.NS{
		Hdr: dns.RR_Header{
			Name:     "newdns.256dpi.com.",
			Rrtype:   dns.TypeNS,
			Class:    dns.ClassINET,
			Ttl:      172800,
			Rdlength: 25,
		},
		Ns: awsNS[2],
	},
	&dns.NS{
		Hdr: dns.RR_Header{
			Name:     "newdns.256dpi.com.",
			Rrtype:   dns.TypeNS,
			Class:    dns.ClassINET,
			Ttl:      172800,
			Rdlength: 22,
		},
		Ns: awsNS[3],
	},
}

var otherNSRRs = []dns.RR{
	&dns.NS{
		Hdr: dns.RR_Header{
			Name:     "other.newdns.256dpi.com.",
			Rrtype:   dns.TypeNS,
			Class:    dns.ClassINET,
			Ttl:      300,
			Rdlength: 23,
		},
		Ns: awsOtherNS[0],
	},
	&dns.NS{
		Hdr: dns.RR_Header{
			Name:     "other.newdns.256dpi.com.",
			Rrtype:   dns.TypeNS,
			Class:    dns.ClassINET,
			Ttl:      300,
			Rdlength: 19,
		},
		Ns: awsOtherNS[1],
	},
	&dns.NS{
		Hdr: dns.RR_Header{
			Name:     "other.newdns.256dpi.com.",
			Rrtype:   dns.TypeNS,
			Class:    dns.ClassINET,
			Ttl:      300,
			Rdlength: 25,
		},
		Ns: awsOtherNS[2],
	},
	&dns.NS{
		Hdr: dns.RR_Header{
			Name:     "other.newdns.256dpi.com.",
			Rrtype:   dns.TypeNS,
			Class:    dns.ClassINET,
			Ttl:      300,
			Rdlength: 22,
		},
		Ns: awsOtherNS[3],
	},
}

func TestAWS(t *testing.T) {
	t.Run("UDP", func(t *testing.T) {
		conformanceTests(t, "udp", awsPrimaryNS+":53", false)
	})

	t.Run("TCP", func(t *testing.T) {
		conformanceTests(t, "tcp", awsPrimaryNS+":53", false)
	})

	t.Run("Resolver", func(t *testing.T) {
		resolverTests(t, awsPrimaryNS+":53")
	})
}

func TestServer(t *testing.T) {
	zone := &Zone{
		Name:             "newdns.256dpi.com.",
		MasterNameServer: awsPrimaryNS,
		AllNameServers: []string{
			awsNS[0],
			awsNS[1],
			awsNS[2],
			awsNS[3],
		},
		AdminEmail: "awsdns-hostmaster@amazon.com",
		Refresh:    2 * time.Hour,
		Retry:      15 * time.Minute,
		Expire:     336 * time.Hour,
		SOATTL:     15 * time.Minute,
		NSTTL:      48 * time.Hour,
		MinTTL:     5 * time.Minute,
		Handler: func(name, remoteAddr string) ([]Set, error) {
			// handle apex records
			if name == "" {
				return []Set{
					{
						Name: "newdns.256dpi.com.",
						Type: A,
						Records: []Record{
							{Address: "1.2.3.4"},
						},
					},
					{
						Name: "newdns.256dpi.com.",
						Type: AAAA,
						Records: []Record{
							{Address: "1:2:3:4::"},
						},
					},
					{
						Name: "newdns.256dpi.com.",
						Type: TXT,
						Records: []Record{
							{Data: []string{"baz"}},
							{Data: []string{"foo", "bar"}},
						},
					},
				}, nil
			}

			// handle example
			if name == "example" {
				return []Set{
					{
						Name: "example.newdns.256dpi.com.",
						Type: CNAME,
						Records: []Record{
							{Address: "example.com."},
						},
					},
				}, nil
			}

			// handle ip4
			if name == "ip4" {
				return []Set{
					{
						Name: "ip4.newdns.256dpi.com.",
						Type: A,
						Records: []Record{
							{Address: "1.2.3.4"},
						},
					},
				}, nil
			}

			// handle ip6
			if name == "ip6" {
				return []Set{
					{
						Name: "ip6.newdns.256dpi.com.",
						Type: AAAA,
						Records: []Record{
							{Address: "1:2:3:4::"},
						},
					},
				}, nil
			}

			// handle mail
			if name == "mail" {
				return []Set{
					{
						Name: "mail.newdns.256dpi.com.",
						Type: MX,
						Records: []Record{
							{Address: "mail.example.com.", Priority: 7},
						},
					},
				}, nil
			}

			// handle multimail
			if name == "multimail" {
				return []Set{
					{
						Name: "multimail.newdns.256dpi.com.",
						Type: MX,
						Records: []Record{
							{Address: "mail1.example.com.", Priority: 1},
							{Address: "mail2.example.com.", Priority: 10},
							{Address: "mail3.example.com.", Priority: 10},
						},
					},
				}, nil
			}

			// handle text
			if name == "text" {
				return []Set{
					{
						Name: "text.newdns.256dpi.com.",
						Type: TXT,
						Records: []Record{
							{Data: []string{"foo", "bar"}},
						},
					},
				}, nil
			}

			// handle ref4
			if name == "ref4" {
				return []Set{
					{
						Name: "ref4.newdns.256dpi.com.",
						Type: CNAME,
						Records: []Record{
							{Address: "ip4.newdns.256dpi.com."},
						},
					},
				}, nil
			}

			// handle ref6
			if name == "ref6" {
				return []Set{
					{
						Name: "ref6.newdns.256dpi.com.",
						Type: CNAME,
						Records: []Record{
							{Address: "ip6.newdns.256dpi.com."},
						},
					},
				}, nil
			}

			// handle refref
			if name == "refref" {
				return []Set{
					{
						Name: "refref.newdns.256dpi.com.",
						Type: CNAME,
						Records: []Record{
							{Address: "ref4.newdns.256dpi.com."},
						},
					},
				}, nil
			}

			// handle ref4m
			if name == "ref4m" {
				return []Set{
					{
						Name: "ref4m.newdns.256dpi.com.",
						Type: MX,
						Records: []Record{
							{Address: "ip4.newdns.256dpi.com.", Priority: 7},
						},
					},
				}, nil
			}

			// handle ref6m
			if name == "ref6m" {
				return []Set{
					{
						Name: "ref6m.newdns.256dpi.com.",
						Type: MX,
						Records: []Record{
							{Address: "ip6.newdns.256dpi.com.", Priority: 7},
						},
					},
				}, nil
			}

			// handle long
			if name == "long" {
				return []Set{
					{
						Name: "long.newdns.256dpi.com.",
						Type: TXT,
						Records: []Record{
							{Data: []string{"gyK4oL9X8Zn3b6TwmUIYAgQx43rBOWMqJWR3wGMGNaZgajnhd2u9JaIbGwNo6gzZunyKYRxID3mKLmYUCcIrNYuo8R4UkijZeshwqEAM2EWnjNsB1hJHOlu6VyRKW13rsFUJedOSqc7YjjUoxm9c3mF28tEXmc3GVsC476wJ2ciSbp7ujDjQ032SQRD6kpayzFX8GncS5KXP8mLK2ZIqK2U4fUmYEpTPQMmp7w24GKkfGJzE4JfMBxSybDUScLq"}},
							{Data: []string{"upNh05zi9flqN2puI9eIGgAgl3gwc65l3WjFdnE3u55dhyUyIoKbOlc1mQJPULPkn1V5TTG9rLBB8AzNfeL8jvwO8h0mzmJhPH8n6dkgI546jB8Z0g0MRJxN5VNSixjFjdR8vtUp6EWlVi7QSe9SYInghV0M17zZ8mXSHwTfYZaPH54ng22mSWzVbRX2tlUPLTNRB5CHrEtxliyhhQlRey98P5G0eo35FUXdqzOSJ3HGqDssBWQAxK3I9feOjbE"}},
							{Data: []string{"z4e6ycRMp6MP3WvWQMxIAOXglxANbj3oB0xD8BffktO4eo3VCR0s6TyGHKixvarOFJU0fqNkXeFOeI7sTXH5X0iXZukfLgnGTxLXNC7KkVFwtVFsh1P0IUNXtNBlOVWrVbxkS62ezbLpENNkiBwbkCvcTjwF2kyI0curAt9JhhJFb3AAq0q1iHWlJLn1KSrev9PIsY3alndDKjYTPxAojxzGKdK3A7rWLJ8Uzb3Z5OhLwP7jTKqbWVUocJRFLYp"}},
						},
					},
				}, nil
			}

			// handle other
			if name == "other" {
				return []Set{
					{
						Name: "other.newdns.256dpi.com.",
						Type: NS,
						Records: []Record{
							{Address: awsOtherNS[0]},
							{Address: awsOtherNS[1]},
							{Address: awsOtherNS[2]},
							{Address: awsOtherNS[3]},
						},
					},
				}, nil
			}

			return nil, nil
		},
	}

	server := NewServer(Config{
		BufferSize: 4096,
		Handler: func(name string) (*Zone, error) {
			if InZone("newdns.256dpi.com.", name) {
				return zone, nil
			}

			return nil, nil
		},
		Logger: func(e Event, msg *dns.Msg, err error, reason string) {
			if e == NetworkError {
				panic(err.Error())
			}
		},
	})

	addr := "0.0.0.0:53001"

	run(server, addr, func() {
		t.Run("UDP", func(t *testing.T) {
			conformanceTests(t, "udp", addr, true)
			additionalTests(t, "udp", addr)
		})

		t.Run("TCP", func(t *testing.T) {
			conformanceTests(t, "tcp", addr, true)
			additionalTests(t, "tcp", addr)
		})

		t.Run("Resolver", func(t *testing.T) {
			resolverTests(t, addr)
		})
	})
}

func TestServerFallback(t *testing.T) {
	zone := &Zone{
		Name:             "example.com.",
		MasterNameServer: "ns1.example.com.",
		AllNameServers: []string{
			"ns1.example.com.",
			"ns2.example.com.",
		},
		Handler: func(name, remoteAddr string) ([]Set, error) {
			// handle apex
			if name == "" {
				return []Set{
					{
						Name: "example.com.",
						Type: A,
						Records: []Record{
							{Address: "1.2.3.4"},
						},
					},
				}, nil
			}

			return nil, nil
		},
	}

	server := NewServer(Config{
		Zones: []string{"example.com."},
		Handler: func(name string) (*Zone, error) {
			if InZone("example.com.", name) {
				return zone, nil
			}

			return nil, nil
		},
		Fallback: "1.1.1.1:53",
	})

	addr := "0.0.0.0:53002"

	run(server, addr, func() {
		// internal zone
		ret, err := Query("udp", addr, "example.com.", "A", func(msg *dns.Msg) {
			msg.RecursionDesired = true
		})
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:         true,
				Authoritative:    true,
				RecursionDesired: true,
			},
			Question: []dns.Question{
				{Name: "example.com.", Qtype: dns.TypeA, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.A{
					Hdr: dns.RR_Header{
						Name:     "example.com.",
						Rrtype:   dns.TypeA,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 4,
					},
					A: net.ParseIP("1.2.3.4"),
				},
			},
			Ns: []dns.RR{
				&dns.NS{
					Hdr: dns.RR_Header{
						Name:     "example.com.",
						Rrtype:   dns.TypeNS,
						Class:    dns.ClassINET,
						Ttl:      172800,
						Rdlength: 6,
					},
					Ns: "ns1.example.com.",
				},
				&dns.NS{
					Hdr: dns.RR_Header{
						Name:     "example.com.",
						Rrtype:   dns.TypeNS,
						Class:    dns.ClassINET,
						Ttl:      172800,
						Rdlength: 6,
					},
					Ns: "ns2.example.com.",
				},
			},
		}, ret)

		// external zone
		ret, err = Query("udp", addr, "newdns.256dpi.com.", "A", func(msg *dns.Msg) {
			msg.RecursionDesired = true
		})
		assert.NoError(t, err)
		ret.Answer[0].(*dns.A).Hdr.Ttl = 300
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:           true,
				Authoritative:      false,
				RecursionDesired:   true,
				RecursionAvailable: true,
			},
			Question: []dns.Question{
				{Name: "newdns.256dpi.com.", Qtype: dns.TypeA, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.A{
					Hdr: dns.RR_Header{
						Name:     "newdns.256dpi.com.",
						Rrtype:   dns.TypeA,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 4,
					},
					A: net.ParseIP("1.2.3.4"),
				},
			},
		}, ret)
	})
}

func conformanceTests(t *testing.T, proto, addr string, local bool) {
	t.Run("ApexA", func(t *testing.T) {
		ret, err := Query(proto, addr, "newdns.256dpi.com.", "A", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "newdns.256dpi.com.", Qtype: dns.TypeA, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.A{
					Hdr: dns.RR_Header{
						Name:     "newdns.256dpi.com.",
						Rrtype:   dns.TypeA,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 4,
					},
					A: net.ParseIP("1.2.3.4"),
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("ApexAAAA", func(t *testing.T) {
		ret, err := Query(proto, addr, "newdns.256dpi.com.", "AAAA", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "newdns.256dpi.com.", Qtype: dns.TypeAAAA, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.AAAA{
					Hdr: dns.RR_Header{
						Name:     "newdns.256dpi.com.",
						Rrtype:   dns.TypeAAAA,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 16,
					},
					AAAA: net.ParseIP("1:2:3:4::"),
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("ApexCNAME", func(t *testing.T) {
		ret, err := Query(proto, addr, "newdns.256dpi.com.", "CNAME", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
				Rcode:         dns.RcodeSuccess,
			},
			Question: []dns.Question{
				{Name: "newdns.256dpi.com.", Qtype: dns.TypeCNAME, Qclass: dns.ClassINET},
			},
			Ns: []dns.RR{
				&dns.SOA{
					Hdr: dns.RR_Header{
						Name:     "newdns.256dpi.com.",
						Rrtype:   dns.TypeSOA,
						Class:    dns.ClassINET,
						Ttl:      900,
						Rdlength: 66,
					},
					Ns:      awsPrimaryNS,
					Mbox:    "awsdns-hostmaster.amazon.com.",
					Serial:  1,
					Refresh: 7200,
					Retry:   900,
					Expire:  1209600,
					Minttl:  300,
				},
			},
		}, ret)
	})

	t.Run("ApexSOA", func(t *testing.T) {
		ret, err := Query(proto, addr, "newdns.256dpi.com.", "SOA", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "newdns.256dpi.com.", Qtype: dns.TypeSOA, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.SOA{
					Hdr: dns.RR_Header{
						Name:     "newdns.256dpi.com.",
						Rrtype:   dns.TypeSOA,
						Class:    dns.ClassINET,
						Ttl:      900,
						Rdlength: 66,
					},
					Ns:      awsPrimaryNS,
					Mbox:    "awsdns-hostmaster.amazon.com.",
					Serial:  1,
					Refresh: 7200,
					Retry:   900,
					Expire:  1209600,
					Minttl:  300,
				},
			},
			Ns: []dns.RR{
				&dns.NS{
					Hdr: dns.RR_Header{
						Name:     "newdns.256dpi.com.",
						Rrtype:   dns.TypeNS,
						Class:    dns.ClassINET,
						Ttl:      172800,
						Rdlength: 23,
					},
					Ns: awsNS[0],
				},
				&dns.NS{
					Hdr: dns.RR_Header{
						Name:     "newdns.256dpi.com.",
						Rrtype:   dns.TypeNS,
						Class:    dns.ClassINET,
						Ttl:      172800,
						Rdlength: 2,
					},
					Ns: awsNS[1],
				},
				&dns.NS{
					Hdr: dns.RR_Header{
						Name:     "newdns.256dpi.com.",
						Rrtype:   dns.TypeNS,
						Class:    dns.ClassINET,
						Ttl:      172800,
						Rdlength: 25,
					},
					Ns: awsNS[2],
				},
				&dns.NS{
					Hdr: dns.RR_Header{
						Name:     "newdns.256dpi.com.",
						Rrtype:   dns.TypeNS,
						Class:    dns.ClassINET,
						Ttl:      172800,
						Rdlength: 22,
					},
					Ns: awsNS[3],
				},
			},
		}, ret)
	})

	t.Run("ApexNS", func(t *testing.T) {
		ret, err := Query(proto, addr, "newdns.256dpi.com.", "NS", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "newdns.256dpi.com.", Qtype: dns.TypeNS, Qclass: dns.ClassINET},
			},
			Answer: nsRRs,
		}, ret)
	})

	t.Run("ApexTXT", func(t *testing.T) {
		ret, err := Query(proto, addr, "newdns.256dpi.com.", "TXT", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "newdns.256dpi.com.", Qtype: dns.TypeTXT, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.TXT{
					Hdr: dns.RR_Header{
						Name:     "newdns.256dpi.com.",
						Rrtype:   dns.TypeTXT,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 4,
					},
					Txt: []string{"baz"},
				},
				&dns.TXT{
					Hdr: dns.RR_Header{
						Name:     "newdns.256dpi.com.",
						Rrtype:   dns.TypeTXT,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 8,
					},
					Txt: []string{"foo", "bar"},
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("SubA", func(t *testing.T) {
		ret, err := Query(proto, addr, "ip4.newdns.256dpi.com.", "A", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "ip4.newdns.256dpi.com.", Qtype: dns.TypeA, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.A{
					Hdr: dns.RR_Header{
						Name:     "ip4.newdns.256dpi.com.",
						Rrtype:   dns.TypeA,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 4,
					},
					A: net.ParseIP("1.2.3.4"),
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("SubAAAA", func(t *testing.T) {
		ret, err := Query(proto, addr, "ip6.newdns.256dpi.com.", "AAAA", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "ip6.newdns.256dpi.com.", Qtype: dns.TypeAAAA, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.AAAA{
					Hdr: dns.RR_Header{
						Name:     "ip6.newdns.256dpi.com.",
						Rrtype:   dns.TypeAAAA,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 16,
					},
					AAAA: net.ParseIP("1:2:3:4::"),
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("SubCNAME", func(t *testing.T) {
		ret, err := Query(proto, addr, "example.newdns.256dpi.com.", "CNAME", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "example.newdns.256dpi.com.", Qtype: dns.TypeCNAME, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.CNAME{
					Hdr: dns.RR_Header{
						Name:     "example.newdns.256dpi.com.",
						Rrtype:   dns.TypeCNAME,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 10,
					},
					Target: "example.com.",
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("SubMX", func(t *testing.T) {
		ret, err := Query(proto, addr, "mail.newdns.256dpi.com.", "MX", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "mail.newdns.256dpi.com.", Qtype: dns.TypeMX, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.MX{
					Hdr: dns.RR_Header{
						Name:     "mail.newdns.256dpi.com.",
						Rrtype:   dns.TypeMX,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 17,
					},
					Mx:         "mail.example.com.",
					Preference: 7,
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("SubMultiMX", func(t *testing.T) {
		ret, err := Query(proto, addr, "multimail.newdns.256dpi.com.", "MX", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "multimail.newdns.256dpi.com.", Qtype: dns.TypeMX, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.MX{
					Hdr: dns.RR_Header{
						Name:     "multimail.newdns.256dpi.com.",
						Rrtype:   dns.TypeMX,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 18,
					},
					Mx:         "mail1.example.com.",
					Preference: 1,
				},
				&dns.MX{
					Hdr: dns.RR_Header{
						Name:     "multimail.newdns.256dpi.com.",
						Rrtype:   dns.TypeMX,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 10,
					},
					Mx:         "mail2.example.com.",
					Preference: 10,
				},
				&dns.MX{
					Hdr: dns.RR_Header{
						Name:     "multimail.newdns.256dpi.com.",
						Rrtype:   dns.TypeMX,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 10,
					},
					Mx:         "mail3.example.com.",
					Preference: 10,
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("SubTXT", func(t *testing.T) {
		ret, err := Query(proto, addr, "text.newdns.256dpi.com.", "TXT", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "text.newdns.256dpi.com.", Qtype: dns.TypeTXT, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.TXT{
					Hdr: dns.RR_Header{
						Name:     "text.newdns.256dpi.com.",
						Rrtype:   dns.TypeTXT,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 8,
					},
					Txt: []string{"foo", "bar"},
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("SubCNAMEForA", func(t *testing.T) {
		ret, err := Query(proto, addr, "example.newdns.256dpi.com.", "A", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "example.newdns.256dpi.com.", Qtype: dns.TypeA, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.CNAME{
					Hdr: dns.RR_Header{
						Name:     "example.newdns.256dpi.com.",
						Rrtype:   dns.TypeCNAME,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 10,
					},
					Target: "example.com.",
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("SubCNAMEForAAAA", func(t *testing.T) {
		ret, err := Query(proto, addr, "example.newdns.256dpi.com.", "AAAA", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "example.newdns.256dpi.com.", Qtype: dns.TypeAAAA, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.CNAME{
					Hdr: dns.RR_Header{
						Name:     "example.newdns.256dpi.com.",
						Rrtype:   dns.TypeCNAME,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 10,
					},
					Target: "example.com.",
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("SubCNAMEForAWithA", func(t *testing.T) {
		ret, err := Query(proto, addr, "ref4.newdns.256dpi.com.", "A", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "ref4.newdns.256dpi.com.", Qtype: dns.TypeA, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.CNAME{
					Hdr: dns.RR_Header{
						Name:     "ref4.newdns.256dpi.com.",
						Rrtype:   dns.TypeCNAME,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 6,
					},
					Target: "ip4.newdns.256dpi.com.",
				},
				&dns.A{
					Hdr: dns.RR_Header{
						Name:     "ip4.newdns.256dpi.com.",
						Rrtype:   dns.TypeA,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 4,
					},
					A: net.ParseIP("1.2.3.4"),
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("SubCNAMEForAAAAWithAAAA", func(t *testing.T) {
		ret, err := Query(proto, addr, "ref6.newdns.256dpi.com.", "AAAA", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "ref6.newdns.256dpi.com.", Qtype: dns.TypeAAAA, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.CNAME{
					Hdr: dns.RR_Header{
						Name:     "ref6.newdns.256dpi.com.",
						Rrtype:   dns.TypeCNAME,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 6,
					},
					Target: "ip6.newdns.256dpi.com.",
				},
				&dns.AAAA{
					Hdr: dns.RR_Header{
						Name:     "ip6.newdns.256dpi.com.",
						Rrtype:   dns.TypeAAAA,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 16,
					},
					AAAA: net.ParseIP("1:2:3:4::"),
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("SubCNAMEWithoutA", func(t *testing.T) {
		ret, err := Query(proto, addr, "ref4.newdns.256dpi.com.", "CNAME", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "ref4.newdns.256dpi.com.", Qtype: dns.TypeCNAME, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.CNAME{
					Hdr: dns.RR_Header{
						Name:     "ref4.newdns.256dpi.com.",
						Rrtype:   dns.TypeCNAME,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 6,
					},
					Target: "ip4.newdns.256dpi.com.",
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("SubCNAMEWithoutAAAA", func(t *testing.T) {
		ret, err := Query(proto, addr, "ref6.newdns.256dpi.com.", "CNAME", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "ref6.newdns.256dpi.com.", Qtype: dns.TypeCNAME, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.CNAME{
					Hdr: dns.RR_Header{
						Name:     "ref6.newdns.256dpi.com.",
						Rrtype:   dns.TypeCNAME,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 6,
					},
					Target: "ip6.newdns.256dpi.com.",
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("SubCNAMEForCNAMEForAWithA", func(t *testing.T) {
		ret, err := Query(proto, addr, "refref.newdns.256dpi.com.", "A", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "refref.newdns.256dpi.com.", Qtype: dns.TypeA, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.CNAME{
					Hdr: dns.RR_Header{
						Name:     "refref.newdns.256dpi.com.",
						Rrtype:   dns.TypeCNAME,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 7,
					},
					Target: "ref4.newdns.256dpi.com.",
				},
				&dns.CNAME{
					Hdr: dns.RR_Header{
						Name:     "ref4.newdns.256dpi.com.",
						Rrtype:   dns.TypeCNAME,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 6,
					},
					Target: "ip4.newdns.256dpi.com.",
				},
				&dns.A{
					Hdr: dns.RR_Header{
						Name:     "ip4.newdns.256dpi.com.",
						Rrtype:   dns.TypeA,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 4,
					},
					A: net.ParseIP("1.2.3.4"),
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("SubMXWithExtraA", func(t *testing.T) {
		ret, err := Query(proto, addr, "ref4m.newdns.256dpi.com.", "MX", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "ref4m.newdns.256dpi.com.", Qtype: dns.TypeMX, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.MX{
					Hdr: dns.RR_Header{
						Name:     "ref4m.newdns.256dpi.com.",
						Rrtype:   dns.TypeMX,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 8,
					},
					Preference: 7,
					Mx:         "ip4.newdns.256dpi.com.",
				},
			},
			Extra: []dns.RR{
				&dns.A{
					Hdr: dns.RR_Header{
						Name:     "ip4.newdns.256dpi.com.",
						Rrtype:   dns.TypeA,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 4,
					},
					A: net.ParseIP("1.2.3.4"),
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("SubMXWithExtraAAAA", func(t *testing.T) {
		ret, err := Query(proto, addr, "ref6m.newdns.256dpi.com.", "MX", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "ref6m.newdns.256dpi.com.", Qtype: dns.TypeMX, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.MX{
					Hdr: dns.RR_Header{
						Name:     "ref6m.newdns.256dpi.com.",
						Rrtype:   dns.TypeMX,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 8,
					},
					Preference: 7,
					Mx:         "ip6.newdns.256dpi.com.",
				},
			},
			Extra: []dns.RR{
				&dns.AAAA{
					Hdr: dns.RR_Header{
						Name:     "ip6.newdns.256dpi.com.",
						Rrtype:   dns.TypeAAAA,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 16,
					},
					AAAA: net.ParseIP("1:2:3:4::"),
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("SubNS", func(t *testing.T) {
		ret, err := Query(proto, addr, "other.newdns.256dpi.com.", "NS", nil)
		assert.NoError(t, err)
		order(ret.Ns)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: false,
			},
			Question: []dns.Question{
				{Name: "other.newdns.256dpi.com.", Qtype: dns.TypeNS, Qclass: dns.ClassINET},
			},
			Ns: order(otherNSRRs),
		}, ret)
	})

	t.Run("NoExactRecord", func(t *testing.T) {
		assertMissing(t, proto, addr, "ip4.newdns.256dpi.com.", "CNAME", dns.RcodeSuccess)
		assertMissing(t, proto, addr, "ip6.newdns.256dpi.com.", "CNAME", dns.RcodeSuccess)
		assertMissing(t, proto, addr, "ip4.newdns.256dpi.com.", "AAAA", dns.RcodeSuccess)
		assertMissing(t, proto, addr, "ip6.newdns.256dpi.com.", "A", dns.RcodeSuccess)
		assertMissing(t, proto, addr, "mail.newdns.256dpi.com.", "A", dns.RcodeSuccess)
		assertMissing(t, proto, addr, "text.newdns.256dpi.com.", "A", dns.RcodeSuccess)
		assertMissing(t, proto, addr, "ip4.newdns.256dpi.com.", "NS", dns.RcodeSuccess)
	})

	t.Run("NoExistingRecord", func(t *testing.T) {
		assertMissing(t, proto, addr, "missing.newdns.256dpi.com.", "A", dns.RcodeNameError)
		assertMissing(t, proto, addr, "missing.newdns.256dpi.com.", "AAAA", dns.RcodeNameError)
		assertMissing(t, proto, addr, "missing.newdns.256dpi.com.", "CNAME", dns.RcodeNameError)
		assertMissing(t, proto, addr, "missing.newdns.256dpi.com.", "MX", dns.RcodeNameError)
		assertMissing(t, proto, addr, "missing.newdns.256dpi.com.", "TXT", dns.RcodeNameError)
		assertMissing(t, proto, addr, "missing.newdns.256dpi.com.", "NS", dns.RcodeNameError)
	})

	t.Run("TruncatedResponse", func(t *testing.T) {
		ret, err := Query(proto, addr, "long.newdns.256dpi.com.", "TXT", nil)
		assert.NoError(t, err)

		if proto == "udp" {
			equalJSON(t, &dns.Msg{
				MsgHdr: dns.MsgHdr{
					Response:      true,
					Authoritative: true,
					Truncated:     true,
				},
				Question: []dns.Question{
					{Name: "long.newdns.256dpi.com.", Qtype: dns.TypeTXT, Qclass: dns.ClassINET},
				},
			}, ret)
		} else {
			equalJSON(t, &dns.Msg{
				MsgHdr: dns.MsgHdr{
					Response:      true,
					Authoritative: true,
				},
				Question: []dns.Question{
					{Name: "long.newdns.256dpi.com.", Qtype: dns.TypeTXT, Qclass: dns.ClassINET},
				},
				Answer: []dns.RR{
					&dns.TXT{
						Hdr: dns.RR_Header{
							Name:     "long.newdns.256dpi.com.",
							Rrtype:   dns.TypeTXT,
							Class:    dns.ClassINET,
							Ttl:      300,
							Rdlength: 256,
						},
						Txt: []string{
							"gyK4oL9X8Zn3b6TwmUIYAgQx43rBOWMqJWR3wGMGNaZgajnhd2u9JaIbGwNo6gzZunyKYRxID3mKLmYUCcIrNYuo8R4UkijZeshwqEAM2EWnjNsB1hJHOlu6VyRKW13rsFUJedOSqc7YjjUoxm9c3mF28tEXmc3GVsC476wJ2ciSbp7ujDjQ032SQRD6kpayzFX8GncS5KXP8mLK2ZIqK2U4fUmYEpTPQMmp7w24GKkfGJzE4JfMBxSybDUScLq",
						},
					},
					&dns.TXT{
						Hdr: dns.RR_Header{
							Name:     "long.newdns.256dpi.com.",
							Rrtype:   dns.TypeTXT,
							Class:    dns.ClassINET,
							Ttl:      300,
							Rdlength: 256,
						},
						Txt: []string{
							"upNh05zi9flqN2puI9eIGgAgl3gwc65l3WjFdnE3u55dhyUyIoKbOlc1mQJPULPkn1V5TTG9rLBB8AzNfeL8jvwO8h0mzmJhPH8n6dkgI546jB8Z0g0MRJxN5VNSixjFjdR8vtUp6EWlVi7QSe9SYInghV0M17zZ8mXSHwTfYZaPH54ng22mSWzVbRX2tlUPLTNRB5CHrEtxliyhhQlRey98P5G0eo35FUXdqzOSJ3HGqDssBWQAxK3I9feOjbE",
						},
					},
					&dns.TXT{
						Hdr: dns.RR_Header{
							Name:     "long.newdns.256dpi.com.",
							Rrtype:   dns.TypeTXT,
							Class:    dns.ClassINET,
							Ttl:      300,
							Rdlength: 256,
						},
						Txt: []string{
							"z4e6ycRMp6MP3WvWQMxIAOXglxANbj3oB0xD8BffktO4eo3VCR0s6TyGHKixvarOFJU0fqNkXeFOeI7sTXH5X0iXZukfLgnGTxLXNC7KkVFwtVFsh1P0IUNXtNBlOVWrVbxkS62ezbLpENNkiBwbkCvcTjwF2kyI0curAt9JhhJFb3AAq0q1iHWlJLn1KSrev9PIsY3alndDKjYTPxAojxzGKdK3A7rWLJ8Uzb3Z5OhLwP7jTKqbWVUocJRFLYp",
						},
					},
				},
				Ns: nsRRs,
			}, ret)
		}
	})

	t.Run("CaseTransfer", func(t *testing.T) {
		ret, err := Query(proto, addr, "Ip4.NeWDnS.256dpi.com.", "A", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "Ip4.NeWDnS.256dpi.com.", Qtype: dns.TypeA, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.A{
					Hdr: dns.RR_Header{
						Name:     "Ip4.NeWDnS.256dpi.com.",
						Rrtype:   dns.TypeA,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 4,
					},
					A: net.ParseIP("1.2.3.4"),
				},
			},
			Ns: []dns.RR{
				&dns.NS{
					Hdr: dns.RR_Header{
						Name:     "NeWDnS.256dpi.com.",
						Rrtype:   dns.TypeNS,
						Class:    dns.ClassINET,
						Ttl:      172800,
						Rdlength: 23,
					},
					Ns: awsNS[0],
				},
				&dns.NS{
					Hdr: dns.RR_Header{
						Name:     "NeWDnS.256dpi.com.",
						Rrtype:   dns.TypeNS,
						Class:    dns.ClassINET,
						Ttl:      172800,
						Rdlength: 19,
					},
					Ns: awsNS[1],
				},
				&dns.NS{
					Hdr: dns.RR_Header{
						Name:     "NeWDnS.256dpi.com.",
						Rrtype:   dns.TypeNS,
						Class:    dns.ClassINET,
						Ttl:      172800,
						Rdlength: 25,
					},
					Ns: awsNS[2],
				},
				&dns.NS{
					Hdr: dns.RR_Header{
						Name:     "NeWDnS.256dpi.com.",
						Rrtype:   dns.TypeNS,
						Class:    dns.ClassINET,
						Ttl:      172800,
						Rdlength: 22,
					},
					Ns: awsNS[3],
				},
			},
		}, ret)
	})

	t.Run("DomainWithSpace", func(t *testing.T) {
		assertMissing(t, proto, addr, "\\ ip4.newdns.256dpi.com.", "NULL", dns.RcodeNameError)
	})

	t.Run("EDNSSuccess", func(t *testing.T) {
		ret, err := Query(proto, addr, "newdns.256dpi.com.", "A", func(msg *dns.Msg) {
			msg.SetEdns0(1337, false)
		})
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "newdns.256dpi.com.", Qtype: dns.TypeA, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.A{
					Hdr: dns.RR_Header{
						Name:     "newdns.256dpi.com.",
						Rrtype:   dns.TypeA,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 4,
					},
					A: net.ParseIP("1.2.3.4"),
				},
			},
			Ns: nsRRs,
			Extra: []dns.RR{
				&dns.OPT{
					Hdr: dns.RR_Header{
						Name:     ".",
						Rrtype:   dns.TypeOPT,
						Class:    4096,
						Ttl:      0,
						Rdlength: 0,
					},
				},
			},
		}, ret)
	})

	t.Run("EDNSError", func(t *testing.T) {
		ret, err := Query(proto, addr, "missing.newdns.256dpi.com.", "A", func(msg *dns.Msg) {
			msg.SetEdns0(1337, false)
		})
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
				Rcode:         dns.RcodeNameError,
			},
			Question: []dns.Question{
				{Name: "missing.newdns.256dpi.com.", Qtype: dns.TypeA, Qclass: dns.ClassINET},
			},
			Ns: []dns.RR{
				&dns.SOA{
					Hdr: dns.RR_Header{
						Name:     "newdns.256dpi.com.",
						Rrtype:   dns.TypeSOA,
						Class:    dns.ClassINET,
						Ttl:      900,
						Rdlength: 66,
					},
					Ns:      awsPrimaryNS,
					Mbox:    "awsdns-hostmaster.amazon.com.",
					Serial:  1,
					Refresh: 7200,
					Retry:   900,
					Expire:  1209600,
					Minttl:  300,
				},
			},
			Extra: []dns.RR{
				&dns.OPT{
					Hdr: dns.RR_Header{
						Name:     ".",
						Rrtype:   dns.TypeOPT,
						Class:    4096,
						Ttl:      0,
						Rdlength: 0,
					},
				},
			},
		}, ret)
	})

	t.Run("EDNSBadVersion", func(t *testing.T) {
		ret, err := Query(proto, addr, "newdns.256dpi.com.", "A", func(msg *dns.Msg) {
			msg.SetEdns0(1337, false)
			msg.Extra[0].(*dns.OPT).SetVersion(2)
		})
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
				Rcode:         dns.RcodeBadVers,
			},
			Question: []dns.Question{
				{Name: "newdns.256dpi.com.", Qtype: dns.TypeA, Qclass: dns.ClassINET},
			},
			Extra: []dns.RR{
				&dns.OPT{
					Hdr: dns.RR_Header{
						Name:     ".",
						Rrtype:   dns.TypeOPT,
						Class:    4096,
						Ttl:      ret.Extra[0].Header().Ttl, // see below
						Rdlength: 0,
					},
				},
			},
		}, ret)

		// the AWS servers sometimes return a bit-flipped TTL value
		ttl := ret.Extra[0].Header().Ttl
		if !local {
			assert.True(t, ttl == 0x1008000 || ttl == 0x1000000, ttl)
		} else {
			assert.Equal(t, uint32(dns.RcodeBadVers<<20), ttl)
		}
	})

	t.Run("EDNSLongResponse", func(t *testing.T) {
		ret, err := Query(proto, addr, "long.newdns.256dpi.com.", "TXT", func(msg *dns.Msg) {
			msg.SetEdns0(1337, false)
		})
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "long.newdns.256dpi.com.", Qtype: dns.TypeTXT, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.TXT{
					Hdr: dns.RR_Header{
						Name:     "long.newdns.256dpi.com.",
						Rrtype:   dns.TypeTXT,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 256,
					},
					Txt: []string{
						"gyK4oL9X8Zn3b6TwmUIYAgQx43rBOWMqJWR3wGMGNaZgajnhd2u9JaIbGwNo6gzZunyKYRxID3mKLmYUCcIrNYuo8R4UkijZeshwqEAM2EWnjNsB1hJHOlu6VyRKW13rsFUJedOSqc7YjjUoxm9c3mF28tEXmc3GVsC476wJ2ciSbp7ujDjQ032SQRD6kpayzFX8GncS5KXP8mLK2ZIqK2U4fUmYEpTPQMmp7w24GKkfGJzE4JfMBxSybDUScLq",
					},
				},
				&dns.TXT{
					Hdr: dns.RR_Header{
						Name:     "long.newdns.256dpi.com.",
						Rrtype:   dns.TypeTXT,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 256,
					},
					Txt: []string{
						"upNh05zi9flqN2puI9eIGgAgl3gwc65l3WjFdnE3u55dhyUyIoKbOlc1mQJPULPkn1V5TTG9rLBB8AzNfeL8jvwO8h0mzmJhPH8n6dkgI546jB8Z0g0MRJxN5VNSixjFjdR8vtUp6EWlVi7QSe9SYInghV0M17zZ8mXSHwTfYZaPH54ng22mSWzVbRX2tlUPLTNRB5CHrEtxliyhhQlRey98P5G0eo35FUXdqzOSJ3HGqDssBWQAxK3I9feOjbE",
					},
				},
				&dns.TXT{
					Hdr: dns.RR_Header{
						Name:     "long.newdns.256dpi.com.",
						Rrtype:   dns.TypeTXT,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 256,
					},
					Txt: []string{
						"z4e6ycRMp6MP3WvWQMxIAOXglxANbj3oB0xD8BffktO4eo3VCR0s6TyGHKixvarOFJU0fqNkXeFOeI7sTXH5X0iXZukfLgnGTxLXNC7KkVFwtVFsh1P0IUNXtNBlOVWrVbxkS62ezbLpENNkiBwbkCvcTjwF2kyI0curAt9JhhJFb3AAq0q1iHWlJLn1KSrev9PIsY3alndDKjYTPxAojxzGKdK3A7rWLJ8Uzb3Z5OhLwP7jTKqbWVUocJRFLYp",
					},
				},
			},
			Ns: nsRRs,
			Extra: []dns.RR{
				&dns.OPT{
					Hdr: dns.RR_Header{
						Name:     ".",
						Rrtype:   dns.TypeOPT,
						Class:    4096,
						Ttl:      0,
						Rdlength: 0,
					},
				},
			},
		}, ret)
	})

	t.Run("RecursionDesired", func(t *testing.T) {
		ret, err := Query(proto, addr, "newdns.256dpi.com.", "A", func(msg *dns.Msg) {
			msg.RecursionDesired = true
		})
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:         true,
				Authoritative:    true,
				RecursionDesired: true,
			},
			Question: []dns.Question{
				{Name: "newdns.256dpi.com.", Qtype: dns.TypeA, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.A{
					Hdr: dns.RR_Header{
						Name:     "newdns.256dpi.com.",
						Rrtype:   dns.TypeA,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 4,
					},
					A: net.ParseIP("1.2.3.4"),
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("UnsupportedMessage", func(t *testing.T) {
		_, err := Query(proto, addr, "newdns.256dpi.com.", "A", func(msg *dns.Msg) {
			msg.Response = true
		})
		assert.True(t, isIOError(err), err)
	})

	t.Run("UnsupportedOpcode", func(t *testing.T) {
		_, err := Query(proto, addr, "newdns.256dpi.com.", "A", func(msg *dns.Msg) {
			msg.Opcode = dns.OpcodeNotify
		})
		assert.True(t, isIOError(err), err)
	})

	t.Run("UnsupportedClass", func(t *testing.T) {
		_, err := Query(proto, addr, "newdns.256dpi.com.", "A", func(msg *dns.Msg) {
			msg.Question[0].Qclass = dns.ClassANY
		})
		assert.True(t, isIOError(err), err)
	})

	t.Run("IgnorePayload", func(t *testing.T) {
		ret, err := Query(proto, addr, "newdns.256dpi.com.", "A", func(msg *dns.Msg) {
			msg.Answer = []dns.RR{
				&dns.A{
					Hdr: dns.RR_Header{
						Name:     "newdns.256dpi.com.",
						Rrtype:   dns.TypeA,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 4,
					},
					A: net.ParseIP("1.2.3.4"),
				},
				&dns.AAAA{
					Hdr: dns.RR_Header{
						Name:     "newdns.256dpi.com.",
						Rrtype:   dns.TypeAAAA,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 4,
					},
					AAAA: net.ParseIP("1:2:3:4::"),
				},
			}
			msg.Ns = []dns.RR{
				nsRRs[0],
			}
			msg.Extra = []dns.RR{
				&dns.TXT{
					Hdr: dns.RR_Header{
						Name:     "newdns.256dpi.com.",
						Rrtype:   dns.TypeTXT,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 4,
					},
					Txt: []string{"baz"},
				},
			}
		})
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
			},
			Question: []dns.Question{
				{Name: "newdns.256dpi.com.", Qtype: dns.TypeA, Qclass: dns.ClassINET},
			},
			Answer: []dns.RR{
				&dns.A{
					Hdr: dns.RR_Header{
						Name:     "newdns.256dpi.com.",
						Rrtype:   dns.TypeA,
						Class:    dns.ClassINET,
						Ttl:      300,
						Rdlength: 4,
					},
					A: net.ParseIP("1.2.3.4"),
				},
			},
			Ns: nsRRs,
		}, ret)
	})

	t.Run("MultipleQuestions", func(t *testing.T) {
		_, err := Query(proto, addr, "newdns.256dpi.com.", "A", func(msg *dns.Msg) {
			msg.Question = append(msg.Question, dns.Question{
				Name:   "newdns.256dpi.com.",
				Qtype:  dns.TypeA,
				Qclass: dns.ClassINET,
			})
		})
		assert.True(t, isIOError(err), err)
	})

	t.Run("UnsupportedType", func(t *testing.T) {
		assertMissing(t, proto, addr, "missing.newdns.256dpi.com.", "NULL", dns.RcodeNameError)
	})

	t.Run("NonAuthoritativeZone", func(t *testing.T) {
		ret, err := Query(proto, addr, "foo.256dpi.com.", "A", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: false,
				Rcode:         dns.RcodeRefused,
			},
			Question: []dns.Question{
				{Name: "foo.256dpi.com.", Qtype: dns.TypeA, Qclass: dns.ClassINET},
			},
		}, ret)
	})
}

func additionalTests(t *testing.T, proto, addr string) {
	t.Run("UnsupportedAny", func(t *testing.T) {
		ret, err := Query(proto, addr, "newdns.256dpi.com.", "ANY", nil)
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:      true,
				Authoritative: true,
				Rcode:         dns.RcodeNotImplemented,
			},
			Question: []dns.Question{
				{Name: "newdns.256dpi.com.", Qtype: dns.TypeANY, Qclass: dns.ClassINET},
			},
		}, ret)
	})
}

func assertMissing(t *testing.T, proto, addr, name, typ string, code int) {
	qt := dns.StringToType[typ]

	ret, err := Query(proto, addr, name, typ, nil)
	assert.NoError(t, err)
	equalJSON(t, &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Response:      true,
			Authoritative: true,
			Rcode:         code,
		},
		Question: []dns.Question{
			{Name: name, Qtype: qt, Qclass: dns.ClassINET},
		},
		Ns: []dns.RR{
			&dns.SOA{
				Hdr: dns.RR_Header{
					Name:     "newdns.256dpi.com.",
					Rrtype:   dns.TypeSOA,
					Class:    dns.ClassINET,
					Ttl:      900,
					Rdlength: 66,
				},
				Ns:      awsPrimaryNS,
				Mbox:    "awsdns-hostmaster.amazon.com.",
				Serial:  1,
				Refresh: 7200,
				Retry:   900,
				Expire:  1209600,
				Minttl:  300,
			},
		},
	}, ret)
}

func resolverTests(t *testing.T, fallback string) {
	var dialer net.Dialer

	resolver := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (conn net.Conn, err error) {
			return dialer.DialContext(ctx, network, fallback)
		},
	}

	ctx := context.Background()

	t.Run("LookupHost", func(t *testing.T) {
		addrs, err := resolver.LookupHost(ctx, "newdns.256dpi.com")
		assert.NoError(t, err)
		assert.Equal(t, []string{"1.2.3.4", "1:2:3:4::"}, addrs)
	})

	t.Run("LookupCNAME", func(t *testing.T) {
		cname, err := resolver.LookupCNAME(ctx, "ref4.newdns.256dpi.com")
		assert.NoError(t, err)
		assert.Equal(t, "ip4.newdns.256dpi.com.", cname)
	})

	t.Run("LookupTXT", func(t *testing.T) {
		txt, err := resolver.LookupTXT(ctx, "newdns.256dpi.com")
		assert.NoError(t, err)
		assert.Equal(t, []string{"baz", "foobar"}, txt)
	})

	t.Run("LookupMX", func(t *testing.T) {
		mx, err := resolver.LookupMX(ctx, "ref4m.newdns.256dpi.com")
		assert.NoError(t, err)
		assert.Equal(t, []*net.MX{
			{Host: "ip4.newdns.256dpi.com.", Pref: 7},
		}, mx)
	})

	t.Run("LookupNS", func(t *testing.T) {
		ns, err := resolver.LookupNS(ctx, "other.newdns.256dpi.com")
		assert.Error(t, err) // zone is not served by server
		assert.Nil(t, ns)
	})
}
