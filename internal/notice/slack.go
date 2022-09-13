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

func (s *Slack) name() string {
	return "Slack"
}

func (s *Slack) buildPayload(r record.Record) string {
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

func (s *Slack) notice(r record.Record) error {
	resp, err := http.Post(s.URL, "application/json", strings.NewReader(s.buildPayload(r)))
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
