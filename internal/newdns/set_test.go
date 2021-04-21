package newdns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetValidate(t *testing.T) {
	table := []struct {
		set Set
		err string
	}{
		{
			set: Set{
				Name: "foo",
			},
			err: "invalid name: foo",
		},
		{
			set: Set{
				Name: "example.com.",
			},
			err: "invalid type: 0",
		},
		{
			set: Set{
				Name: "example.com.",
				Type: A,
			},
			err: "missing records",
		},
		{
			set: Set{
				Name: "example.com.",
				Type: A,
				Records: []Record{
					{Address: "foo"},
				},
			},
			err: "invalid record: invalid IPv4 address: foo",
		},
		{
			set: Set{
				Name: "example.com.",
				Type: TXT,
				Records: []Record{
					{},
				},
			},
			err: "invalid record: missing data",
		},
		{
			set: Set{
				Name: "example.com.",
				Type: CNAME,
				Records: []Record{
					{},
					{},
				},
			},
			err: "multiple CNAME records",
		},
		{
			set: Set{
				Name: "example.com.",
				Type: A,
				Records: []Record{
					{Address: "1.2.3.4"},
					{Address: "1.2.3.4"},
				},
			},
			err: "duplicate address: 1.2.3.4",
		},
	}

	for i, item := range table {
		err := item.set.Validate()
		if err != nil {
			assert.Equal(t, item.err, err.Error(), i)
		} else {
			assert.Equal(t, item.err, "", item)
		}
	}
}
