package newdns

import (
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func TestResolver(t *testing.T) {
	ret, err := Query("tcp", "1.1.1.1:53", "example.newdns.256dpi.com.", "A", func(msg *dns.Msg) {
		msg.RecursionDesired = true
	})
	assert.NoError(t, err)
	equalJSON(t, &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Response:           true,
			RecursionDesired:   true,
			RecursionAvailable: true,
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
					Ttl:      ret.Answer[0].(*dns.CNAME).Hdr.Ttl,
					Rdlength: 10,
				},
				Target: "example.com.",
			},
			&dns.A{
				Hdr: dns.RR_Header{
					Name:     "example.com.",
					Rrtype:   dns.TypeA,
					Class:    dns.ClassINET,
					Ttl:      ret.Answer[1].(*dns.A).Hdr.Ttl,
					Rdlength: 4,
				},
				A: ret.Answer[1].(*dns.A).A,
			},
		},
	}, ret)

	addr := "0.0.0.0:53002"
	mux := dns.NewServeMux()
	mux.Handle("newdns.256dpi.com", Proxy(awsNS[0]+":53", nil))
	mux.Handle("example.com", Proxy("a.iana-servers.net:53", nil))
	handler := Resolver(mux)

	serve(handler, addr, func() {
		ret, err := Query("udp", addr, "example.newdns.256dpi.com.", "A", func(msg *dns.Msg) {
			msg.RecursionDesired = true
		})
		assert.NoError(t, err)
		equalJSON(t, &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Response:           true,
				RecursionDesired:   true,
				RecursionAvailable: true,
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
						Ttl:      ret.Answer[0].(*dns.CNAME).Hdr.Ttl,
						Rdlength: 13,
					},
					Target: "example.com.",
				},
				&dns.A{
					Hdr: dns.RR_Header{
						Name:     "example.com.",
						Rrtype:   dns.TypeA,
						Class:    dns.ClassINET,
						Ttl:      ret.Answer[1].(*dns.A).Hdr.Ttl,
						Rdlength: 4,
					},
					A: ret.Answer[1].(*dns.A).A,
				},
			},
		}, ret)
	})
}
