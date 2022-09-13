package notice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/li4n0/revsuit/internal/record"
	"github.com/pkg/errors"
)

var _ Bot = (*Lark)(nil)

type Lark struct {
	URL string
}

type larkText struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type larkElement struct {
	Tag  string   `json:"tag"`
	Text larkText `json:"text"`
}

type larkCard struct {
	Header   larkHeader    `json:"header"`
	Elements []larkElement `json:"elements"`
}

type larkHeader struct {
	Title larkText `json:"title"`
}

type larkPayload struct {
	MsgType string   `json:"msg_type"`
	Card    larkCard `json:"card"`
}

func (l *Lark) name() string {
	return "Lark"
}

func (l *Lark) buildPayload(r record.Record) string {
	payload := larkPayload{
		MsgType: "interactive",
		Card: larkCard{
			Header: larkHeader{
				Title: larkText{
					Tag:     "plain_text",
					Content: "New Connection",
				},
			},
			Elements: []larkElement{
				{
					Tag: "div",
					Text: larkText{
						Tag:     "lark_md",
						Content: formatRecordField(r, "**%s**: %v"),
					},
				},
			},
		},
	}
	p, err := json.Marshal(&payload)
	if err != nil {
		return ""
	}
	return string(p)
}

func (l *Lark) notice(r record.Record) error {
	resp, err := http.Post(l.URL, "application/json", strings.NewReader(l.buildPayload(r)))
	if err != nil {
		return errors.Wrap(err, "HTTP request")
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "read HTTP response body")
		}
		return fmt.Errorf("non-success response status code %d with body: %s", resp.StatusCode, data)
	}
	return nil
}
