package newdns

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func run(s *Server, addr string, fn func()) {
	defer s.Close()

	go func() {
		err := s.Run(addr)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	fn()
}

func serve(handler dns.Handler, addr string, fn func()) {
	closer := make(chan struct{})
	defer close(closer)

	go func() {
		err := Run(addr, handler, Accept(nil), closer)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	fn()
}

func equalJSON(t *testing.T, a, b interface{}) {
	buf := new(bytes.Buffer)

	e := json.NewEncoder(buf)
	e.SetIndent("", "  ")

	_ = e.Encode(a)
	aa := buf.String()

	buf.Reset()
	_ = e.Encode(b)
	bb := buf.String()

	assert.JSONEq(t, aa, bb)
}

func order(rrs []dns.RR) []dns.RR {
	cpy := make([]dns.RR, len(rrs))
	copy(cpy, rrs)

	sort.Slice(cpy, func(i, j int) bool {
		var ai string
		switch rr := cpy[i].(type) {
		case *dns.NS:
			ai = rr.Ns
			rr.Hdr.Rdlength = 0
		}

		var aj string
		switch rr := cpy[j].(type) {
		case *dns.NS:
			aj = rr.Ns
			rr.Hdr.Rdlength = 0
		}

		return ai < aj
	})

	return cpy
}

func isIOError(err error) bool {
	if err == nil {
		return false
	}

	if strings.Contains(err.Error(), "i/o timeout") {
		return true
	}

	if strings.Contains(err.Error(), "connection reset by peer") {
		return true
	}

	return false
}
