package newdns

import "github.com/miekg/dns"

// Proxy returns a handler that proxies requests to the provided DNS server. The
// optional logger is called with events about the processing of requests.
func Proxy(addr string, logger Logger) dns.Handler {
	return dns.HandlerFunc(func(w dns.ResponseWriter, req *dns.Msg) {
		// log request
		if logger != nil {
			logger(ProxyRequest, req, nil, "")
		}

		// forward request to fallback
		rs, err := dns.Exchange(req, addr)
		if err != nil {
			if logger != nil {
				logger(ProxyError, nil, err, "")
			}
			_ = w.Close()
			return
		}

		// log response
		if logger != nil {
			logger(ProxyResponse, rs, nil, "")
		}

		// write response
		err = w.WriteMsg(rs)
		if err != nil {
			if logger != nil {
				logger(NetworkError, nil, err, "")
			}
			_ = w.Close()
		}
	})
}
