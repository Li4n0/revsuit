package newdns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecordValidate(t *testing.T) {
	table := []struct {
		typ Type
		rec Record
		err string
	}{
		{
			typ: A,
			rec: Record{Address: "foo"},
			err: "invalid IPv4 address: foo",
		},
		{
			typ: AAAA,
			rec: Record{Address: "foo"},
			err: "invalid IPv6 address: foo",
		},
		{
			typ: A,
			rec: Record{Address: "1:2:3:4::"},
			err: "invalid IPv4 address: 1:2:3:4::",
		},
		{
			typ: A,
			rec: Record{Address: "1.2.3.4"},
		},
		{
			typ: AAAA,
			rec: Record{Address: "1:2:3:4::"},
		},
		{
			typ: CNAME,
			rec: Record{Address: "---"},
			err: "invalid domain name: ---",
		},
		{
			typ: CNAME,
			rec: Record{Address: "foo.com"},
			err: "invalid domain name: foo.com",
		},
		{
			typ: CNAME,
			rec: Record{Address: "foo.com."},
		},
		{
			typ: MX,
			rec: Record{Address: "foo.com"},
			err: "invalid domain name: foo.com",
		},
		{
			typ: MX,
			rec: Record{Address: "foo.com."},
		},
		{
			typ: TXT,
			rec: Record{Data: nil},
			err: "missing data",
		},
		{
			typ: TXT,
			rec: Record{Data: []string{"z4e6ycRMp6MP3WvWQMxIAOXglxANbj3oB0xD8BffktO4eo3VCR0s6TyGHKixvarOFJU0fqNkXeFOeI7sTXH5X0iXZukfLgnGTxLXNC7KkVFwtVFsh1P0IUNXtNBlOVWrVbxkS62ezbLpENNkiBwbkCvcTjwF2kyI0curAt9JhhJFb3AAq0q1iHWlJLn1KSrev9PIsY3alndDKjYTPxAojxzGKdK3A7rWLJ8Uzb3Z5OhLwP7jTKqbWVUocJRFLYpL"}},
			err: "data too long",
		},
		{
			typ: TXT,
			rec: Record{Data: []string{"foo"}},
		},
		{
			typ: NS,
			rec: Record{Address: "foo.com"},
			err: "invalid ns name: foo.com",
		},
		{
			typ: NS,
			rec: Record{Address: "foo.com."},
		},
	}

	for i, item := range table {
		err := item.rec.Validate(item.typ)
		if err != nil {
			assert.Equal(t, item.err, err.Error(), i)
		} else {
			assert.Equal(t, item.err, "", item)
		}
	}
}
