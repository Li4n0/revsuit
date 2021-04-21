package newdns

import (
	"time"

	"github.com/miekg/dns"
)

// Query can be used to query a DNS server over the provided protocol on its
// address for the specified name and type. The supplied function can be set to
// mutate the sent request.
func Query(proto, addr, name, typ string, fn func(*dns.Msg)) (*dns.Msg, error) {
	// prepare request
	req := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Id: dns.Id(),
		},
		Question: []dns.Question{
			{
				Name:   name,
				Qtype:  dns.StringToType[typ],
				Qclass: dns.ClassINET,
			},
		},
	}

	// call function if available
	if fn != nil {
		fn(req)
	}

	// prepare client
	client := dns.Client{
		Net:     proto,
		Timeout: time.Second,
	}

	// send request
	res, _, err := client.Exchange(req, addr)
	if err != nil {
		return nil, err
	}

	// reset id to allow direct comparison
	res.Id = 0

	return res, nil
}
