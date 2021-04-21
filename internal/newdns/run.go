package newdns

import (
	"fmt"

	"github.com/miekg/dns"
)

// Accept will return a dns.MsgAcceptFunc that only accepts normal queries.
func Accept(logger Logger) dns.MsgAcceptFunc {
	return func(dh dns.Header) dns.MsgAcceptAction {
		// check if request
		if dh.Bits&(1<<15) != 0 {
			log(logger, Ignored, nil, nil, fmt.Sprintf("not a request"))
			return dns.MsgIgnore
		}

		// check opcode
		if int(dh.Bits>>11)&0xF != dns.OpcodeQuery {
			log(logger, Ignored, nil, nil, fmt.Sprintf("not a query"))
			return dns.MsgIgnore
		}

		// check question count
		if dh.Qdcount != 1 {
			log(logger, Ignored, nil, nil, fmt.Sprintf("invalid question count: %d", dh.Qdcount))
			return dns.MsgIgnore
		}

		return dns.MsgAccept
	}
}

// Run will start a UDP and TCP listener to serve the specified handler with the
// specified accept function until the provided close channel is closed. It will
// return the first error of a listener.
func Run(addr string, handler dns.Handler, accept dns.MsgAcceptFunc, close <-chan struct{}) error {
	// prepare servers
	udp := &dns.Server{Addr: addr, Net: "udp", Handler: handler, MsgAcceptFunc: accept}
	tcp := &dns.Server{Addr: addr, Net: "tcp", Handler: handler, MsgAcceptFunc: accept}

	// prepare errors
	errs := make(chan error, 2)

	// run udp server
	go func() {
		errs <- udp.ListenAndServe()
	}()

	// run tcp server
	go func() {
		errs <- tcp.ListenAndServe()
	}()

	// await first error
	var err error
	select {
	case err = <-errs:
	case <-close:
	}

	// shutdown servers
	_ = udp.Shutdown()
	_ = tcp.Shutdown()

	return err
}
