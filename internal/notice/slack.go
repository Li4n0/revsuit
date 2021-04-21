package notice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/li4n0/revsuit/internal/record"
)

var _ Bot = (*Slack)(nil)

type Slack struct {
	URL string
}

type slackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type slackBlock struct {
	Type string    `json:"type"`
	Text slackText `json:"text"`
}

type slackAttachments struct {
	Color  string       `json:"color"`
	Blocks []slackBlock `json:"blocks"`
}

type slackPayload struct {
	Attachments []slackAttachments `json:"attachments"`
}

func (d *Slack) buildPayload(r record.Record) string {
	payload := slackPayload{
		Attachments: []slackAttachments{
			{
				Color: "#f2c744",
				Blocks: []slackBlock{
					{
						Type: "header",
						Text: slackText{
							Type: "plain_text",
							Text: "New Connection",
						},
					},
					{
						Type: "section",
						Text: slackText{
							Type: "mrkdwn",
							Text: formatRecordField(r, "- `%s: %v`"),
						},
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

func (d *Slack) notice(r record.Record) error {
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
