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

var _ Bot = (*DingTalk)(nil)

type DingTalk struct {
	URL string
}

type dingAt struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type dingMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type dingPayload struct {
	MsgType  string       `json:"msgtype"`
	Markdown dingMarkdown `json:"markdown"`
	At       []dingAt     `json:"at"`
}

func (d *DingTalk) name() string {
	return "DingTalk"
}

func (d *DingTalk) buildPayload(r record.Record) string {
	payload := dingPayload{
		MsgType: "markdown",
		Markdown: dingMarkdown{
			Title: "New Connection",
			Text: "**<font color='#e96900' face='Fira Code' size='3'>New Connection</font>**\n" +
				formatRecordField(r, "> **<font color='#e96900' face='Fira Code'>%s: </font>**<font color='#e96900' face='Fira Code'>%v</font>\n"),
		},
		At: []dingAt{
			{
				IsAtAll: true,
			},
		},
	}
	p, err := json.Marshal(&payload)
	if err != nil {
		return ""
	}
	return string(p)
}

func (d *DingTalk) notice(r record.Record) error {
	resp, err := http.Post(d.URL, "application/json", strings.NewReader(d.buildPayload(r)))
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
