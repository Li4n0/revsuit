# newdns

[![Build Status](https://travis-ci.org/256dpi/newdns.svg?branch=master)](https://travis-ci.org/256dpi/newdns)
[![Coverage Status](https://coveralls.io/repos/github/256dpi/newdns/badge.svg?branch=master)](https://coveralls.io/github/256dpi/newdns?branch=master)
[![GoDoc](https://godoc.org/github.com/256dpi/newdns?status.svg)](http://godoc.org/github.com/256dpi/newdns)
[![Release](https://img.shields.io/github/release/256dpi/newdns.svg)](https://github.com/256dpi/newdns/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/256dpi/newdns)](https://goreportcard.com/report/github.com/256dpi/newdns)
 
**A library for building custom DNS servers in Go.**

The newdns library wraps the widely used, but low-level [github.com/miekg/dns](https://github.com/miekg/dns) package with a simple interface to quickly build custom DNS servers. The implemented server only supports a subset of record types (A, AAAA, CNAME, MX, TXT) and is intended to be used as a leaf authoritative name server only. It supports UDP and TCP as transport protocols and implements EDNS0. Conformance is tested by issuing a corpus of tests against a zone in AWS Route53 and comparing the response and behavior.

The intention of this project is not to build a feature-complete alternative to "managed zone" offerings by major cloud platforms. However, some projects may require frequent synchronization of many records between a custom database and a cloud-hosted "managed zone". In this scenario, a custom DNS server that queries the own database might be a lot simpler to manage and operate. Also, the distributed nature of the DNS system offers interesting qualities that could be leveraged by future applications.

## Example

```go
// create zone
zone := &newdns.Zone{
    Name:             "example.com.",
    MasterNameServer: "ns1.hostmaster.com.",
    AllNameServers: []string{
        "ns1.hostmaster.com.",
        "ns2.hostmaster.com.",
        "ns3.hostmaster.com.",
    },
    Handler: func(name string) ([]newdns.Set, error) {
        // return apex records
        if name == "" {
            return []newdns.Set{
                {
                    Name: "example.com.",
                    Type: newdns.A,
                    Records: []newdns.Record{
                        {Address: "1.2.3.4"},
                    },
                },
                {
                    Name: "example.com.",
                    Type: newdns.AAAA,
                    Records: []newdns.Record{
                        {Address: "1:2:3:4::"},
                    },
                },
            }, nil
        }

        // return sub records
        if name == "foo" {
            return []newdns.Set{
                {
                    Name: "foo.example.com.",
                    Type: newdns.CNAME,
                    Records: []newdns.Record{
                        {Address: "bar.example.com."},
                    },
                },
            }, nil
        }

        return nil, nil
    },
}

// create server
server := newdns.NewServer(newdns.Config{
    Handler: func(name string) (*newdns.Zone, error) {
        // check name
        if newdns.InZone("example.com.", name) {
            return zone, nil
        }

        return nil, nil
    },
    Logger: func(e newdns.Event, msg *dns.Msg, err error, reason string) {
        fmt.Println(e, err, reason)
    },
})

// run server
go func() {
    err := server.Run(":1337")
    if err != nil {
        panic(err)
    }
}()

// print info
fmt.Println("Query apex: dig example.com @0.0.0.0 -p 1337")
fmt.Println("Query other: dig foo.example.com @0.0.0.0 -p 1337")

// wait forever
select {}
```

## Credits

- https://github.com/miekg/dns
- https://github.com/coredns/coredns
