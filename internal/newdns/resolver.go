package newdns

import (
	"net"

	"github.com/miekg/dns"
)

var fakeAddr = &net.TCPAddr{
	IP:   net.IP{0, 0, 0, 0},
	Port: 0,
}

type responseWriter struct {
	msg *dns.Msg
}

func (w *responseWriter) LocalAddr() net.Addr {
	return fakeAddr
}

func (w *responseWriter) RemoteAddr() net.Addr {
	return fakeAddr
}

func (w *responseWriter) WriteMsg(msg *dns.Msg) error {
	// check message
	if w.msg != nil {
		panic("message already set")
	}

	// set message
	w.msg = msg

	return nil
}

func (w *responseWriter) Write([]byte) (int, error) {
	panic("not implemented")
}

func (w *responseWriter) Close() error {
	return nil
}

func (w *responseWriter) TsigStatus() error {
	panic("not implemented")
}

func (w *responseWriter) TsigTimersOnly(bool) {
	panic("not implemented")
}

func (w *responseWriter) Hijack() {
	panic("not implemented")
}

// Resolver returns a very primitive recursive resolver that uses the provided
// handler to resolve all names.
func Resolver(handler dns.Handler) dns.Handler {
	return dns.HandlerFunc(func(w dns.ResponseWriter, req *dns.Msg) {
		// forward query if no recursion is desired
		if !req.RecursionDesired {
			handler.ServeDNS(w, req)
			return
		}

		// prepare response
		res := new(dns.Msg)
		res.SetReply(req)
		res.RecursionAvailable = true

		// query handler
		var wr responseWriter
		handler.ServeDNS(&wr, req)

		// check response
		if wr.msg == nil {
			_ = w.WriteMsg(res)
			return
		}

		// add resolved answers
		res.Answer = append(res.Answer, resolve(handler, wr.msg.Answer)...)

		// write response
		err := w.WriteMsg(res)
		if err != nil {
			_ = w.Close()
		}
	})
}

func resolve(handler dns.Handler, records []dns.RR) []dns.RR {
	// prepare result
	var res []dns.RR
	res = append(res, records...)

	// handle records
	for _, record := range records {
		if cname, ok := record.(*dns.CNAME); ok {
			// query handler
			var wr responseWriter
			handler.ServeDNS(&wr, &dns.Msg{
				Question: []dns.Question{
					{
						Name:   cname.Target,
						Qtype:  dns.TypeA,
						Qclass: dns.ClassINET,
					},
				},
			})

			// add resolved answers
			res = append(res, resolve(handler, wr.msg.Answer)...)
		}
	}

	return res
}
