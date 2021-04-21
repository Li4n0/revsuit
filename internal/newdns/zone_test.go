package newdns

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZoneValidate(t *testing.T) {
	table := []struct {
		zne Zone
		err string
	}{
		{
			zne: Zone{
				Name: "foo",
			},
			err: "name not fully qualified: foo",
		},
		{
			zne: Zone{
				Name:             "example.com.",
				MasterNameServer: "foo",
			},
			err: "master server not full qualified: foo",
		},
		{
			zne: Zone{
				Name:             "example.com.",
				MasterNameServer: "n1.example.com.",
			},
			err: "missing name servers",
		},
		{
			zne: Zone{
				Name:             "example.com.",
				MasterNameServer: "n1.example.com.",
				AllNameServers: []string{
					"foo",
				},
			},
			err: "name server not fully qualified: foo",
		},
		{
			zne: Zone{
				Name:             "example.com.",
				MasterNameServer: "n2.example.com.",
				AllNameServers: []string{
					"n1.example.com.",
				},
			},
			err: "master name server not listed as name server: n2.example.com.",
		},
		{
			zne: Zone{
				Name:             "example.com.",
				MasterNameServer: "n1.example.com.",
				AllNameServers: []string{
					"n1.example.com.",
				},
			},
		},
		{
			zne: Zone{
				Name:             "example.com.",
				MasterNameServer: "n1.example.com.",
				AllNameServers: []string{
					"n1.example.com.",
				},
				AdminEmail: "foo@bar..example.com",
			},
			err: "admin email cannot be converted to a domain name: foo@bar..example.com",
		},
		{
			zne: Zone{
				Name:             "example.com.",
				MasterNameServer: "n1.example.com.",
				AllNameServers: []string{
					"n1.example.com.",
				},
				Refresh: 1,
				Retry:   2,
			},
			err: "retry must be less than refresh: 2",
		},
		{
			zne: Zone{
				Name:             "example.com.",
				MasterNameServer: "n1.example.com.",
				AllNameServers: []string{
					"n1.example.com.",
				},
				Expire: 1,
				Retry:  2,
			},
			err: "expire must be bigger than the sum of refresh and retry: 1",
		},
	}

	for i, item := range table {
		err := item.zne.Validate()
		if err != nil {
			assert.EqualValues(t, item.err, err.Error(), i)
		} else {
			assert.Equal(t, item.err, "", item)
		}
	}
}

func TestZoneLookup(t *testing.T) {
	zone := Zone{
		Name:             "example.com.",
		MasterNameServer: "ns1.example.com.",
		AllNameServers: []string{
			"ns1.example.com.",
			"ns2.example.com.",
		},
		Handler: func(name, remoteAddr string) ([]Set, error) {
			if name == "error" {
				return nil, io.EOF
			}

			if name == "invalid1" {
				return []Set{
					{Name: "foo"},
				}, nil
			}

			if name == "invalid2" {
				return []Set{
					{Name: "foo.", Type: A, Records: []Record{{Address: "1.2.3.4"}}},
				}, nil
			}

			if name == "multiple" {
				return []Set{
					{Name: "foo.example.com.", Type: A, Records: []Record{{Address: "1.2.3.4"}}},
					{Name: "foo.example.com.", Type: A, Records: []Record{{Address: "1.2.3.4"}}},
				}, nil
			}

			if name == "" {
				return []Set{
					{Name: "example.com.", Type: CNAME, Records: []Record{{Address: "cool.com."}}},
				}, nil
			}

			if name == "cname" {
				return []Set{
					{Name: "cname.example.com.", Type: A, Records: []Record{{Address: "1.2.3.4"}}},
					{Name: "cname.example.com.", Type: CNAME, Records: []Record{{Address: "cool.com."}}},
				}, nil
			}

			return nil, nil
		},
	}

	err := zone.Validate()
	assert.NoError(t, err)

	table := []struct {
		name string
		err  string
	}{
		{
			name: "foo",
			err:  "invalid name: foo",
		},
		{
			name: "foo.",
			err:  "name does not belong to zone: foo.",
		},
		{
			name: "error.example.com.",
			err:  "zone handler error: EOF",
		},
		{
			name: "invalid1.example.com.",
			err:  "invalid set: invalid name: foo",
		},
		{
			name: "invalid2.example.com.",
			err:  "set does not belong to zone: foo.",
		},
		{
			name: "multiple.example.com.",
			err:  "multiple sets for same type",
		},
		{
			name: "example.com.",
			err:  "invalid CNAME set at apex: example.com.",
		},
		{
			name: "cname.example.com.",
			err:  "other sets with CNAME set: cname.example.com.",
		},
	}

	for i, item := range table {
		res, exists, err := zone.Lookup(item.name, "127.0.0.1", A)
		assert.Equal(t, item.err, err.Error(), i)
		assert.False(t, exists, i)
		assert.Nil(t, res, i)
	}
}
