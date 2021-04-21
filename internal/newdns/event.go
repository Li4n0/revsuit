package newdns

import "github.com/miekg/dns"

// Event denotes an event type emitted to the logger.
type Event int

const (
	// Ignored are requests that haven been dropped by leaving the connection
	// hanging to mitigate attacks. Inspect the reason for more information.
	Ignored Event = iota

	// Request is emitted for every accepted request. For every request event
	// a finish event fill follow. You can inspect the message to see the
	// complete request sent by the client.
	Request Event = iota

	// Refused are requests that received an error due to some incompatibility.
	// Inspect the reason for more information.
	Refused Event = iota

	// BackendError is emitted with errors returned by the callback and
	// validation functions. Inspect the error for more information.
	BackendError Event = iota

	// NetworkError is emitted with errors returned by the connection. Inspect
	// the error for more information.
	NetworkError Event = iota

	// Response is emitted with the final response to the client. You can inspect
	// the message to see the complete response to the client.
	Response Event = iota

	// Finish is emitted when a request has been processed.
	Finish Event = iota

	// ProxyRequest is emitted with every request forwarded to the fallback
	// DNS server.
	ProxyRequest Event = iota

	// ProxyResponse is emitted with ever response received from the fallback
	// DNS server.
	ProxyResponse Event = iota

	// ProxyError is emitted with errors returned by the fallback DNS server.
	// Inspect the error for more information.
	ProxyError Event = iota
)

// String will return the name of the event.
func (e Event) String() string {
	switch e {
	case Ignored:
		return "Ignored"
	case Request:
		return "Request"
	case Refused:
		return "Refused"
	case BackendError:
		return "BackendError"
	case NetworkError:
		return "NetworkError"
	case Response:
		return "Response"
	case Finish:
		return "Finish"
	case ProxyRequest:
		return "ProxyRequest"
	case ProxyResponse:
		return "ProxyResponse"
	case ProxyError:
		return "ProxyError"
	default:
		return "Unknown"
	}
}

// Logger is function that accepts logging events.
type Logger func(e Event, msg *dns.Msg, err error, reason string)

func log(l Logger, e Event, msg *dns.Msg, err error, reason string) {
	if l != nil {
		l(e, msg, err, reason)
	}
}
