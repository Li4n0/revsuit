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

var _ Bot = (*Weixin)(nil)

type Weixin struct {
	URL string
}

type weixinMarkdown struct {
	Content string `json:"content"`
}

type weixinPayload struct {
	ToUser   string         `json:"touser"`
	MsgType  string         `json:"msgtype"`
	Markdown weixinMarkdown `json:"markdown"`
}

func (w *Weixin) name() string {
	return "Weixin"
}

func (w *Weixin) buildPayload(r record.Record) string {
	payload := weixinPayload{
		ToUser:  "@all",
		MsgType: "markdown",
		Markdown: weixinMarkdown{
			Content: "<font color='#e96900' face='Fira Code' size=3>New Connection</font>\n" +
				formatRecordField(r, `> **<font color="#e96900" face="Fira Code">%s: </font>**<font color="#e96900" face="Fira Code">%v</font>`),
		},
	}
	p, err := json.Marshal(&payload)
	if err != nil {
		return ""
	}
	return string(p)
}

func (w *Weixin) notice(r record.Record) error {
	resp, err := http.Post(w.URL, "application/json", strings.NewReader(w.buildPayload(r)))
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
