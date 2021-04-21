package notice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/li4n0/revsuit/internal/record"
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

func (d *Lark) buildPayload(r record.Record) string {
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

func (d *Lark) notice(r record.Record) error {
	resp, err := http.DefaultClient.Post(d.URL, "application/json", strings.NewReader(d.buildPayload(r)))
	if err != nil {
		return fmt.Errorf("HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read HTTP response body: %v", err)
		}
		return fmt.Errorf("non-success response status code %d with body: %s", resp.StatusCode, data)
	}
	return nil
}
