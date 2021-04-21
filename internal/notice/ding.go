package notice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/li4n0/revsuit/internal/record"
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

func (d *DingTalk) buildPayload(r record.Record) string {
	payload := dingPayload{
		MsgType: "markdown",
		Markdown: dingMarkdown{
			Title: "New Connection",
			Text: "**<font color=\"#e96900\" face=\"Fira Code\" size=\"3\">New Connection</font>**\n" +
				formatRecordField(r, "> **<font color=\"#e96900\" face=\"Fira Code\">%s: </font>**<font color=\"#e96900\" face=\"Fira Code\">%v</font>\n"),
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
