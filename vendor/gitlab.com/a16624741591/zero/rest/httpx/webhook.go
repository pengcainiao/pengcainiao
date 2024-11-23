package httpx

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"gitlab.com/a16624741591/zero/core/env"
	"gitlab.com/a16624741591/zero/core/logx"
	"gitlab.com/a16624741591/zero/core/sysx"
)

// WebHookRebot 企业微信webhook 机器人
type WebHookRebot string

// MessageType 消息类型
type MessageType string
type ReportErrors struct {
	Title   string                 `json:"title,omitempty"`   //标题
	Payload string                 `json:"payload,omitempty"` //内容
	Args    interface{}            `json:"args,omitempty"`    //参数
	Error   error                  `json:"error,omitempty"`   //错误
	Extra   map[string]interface{} `json:"extra,omitempty"`   //其他相关数据
}

type Message struct {
	MsgType MessageType `json:"msgtype,omitempty"`
	Text    struct {
		Content string `json:"content,omitempty"`
	} `json:"text,omitempty"`
	Markdown struct {
		Content string `json:"content,omitempty"`
	} `json:"markdown,omitempty"`
}

const (
	TxtMessage      MessageType = "text"
	MarkdownMessage MessageType = "markdown"
)

const (
	develop          WebHookRebot = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=ca4e7f42-d76b-49dc-bbb4-1b8380b9068a"
	production       WebHookRebot = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=52b9f758-f128-4a63-bdad-24dc95013871"
	CustomerFeedback WebHookRebot = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=e20e76ec-f5c6-4387-9964-0d1f204593c8"
)

func WebHook() WebHookRebot {
	var env = os.Getenv("RELEASE_MODE")
	if env == "production" {
		return production
	}
	return develop
}

func (m ReportErrors) String() string {
	return m.writeString(false)
}

func (m ReportErrors) MarkdownString() string {
	return m.writeString(true)
}

func (m ReportErrors) writeString(markdown bool) string {
	var (
		buf   = &bytes.Buffer{}
		tag   string
		title string
	)
	if markdown {
		tag = "> "
		title = "#### "
	}
	if m.Title != "" {
		buf.WriteString(fmt.Sprintf("%s%s \n", title, m.Title))
	}
	buf.WriteString(fmt.Sprintf("%s[服务]%s(%s) \n", tag, sysx.SubSystem, env.ReleaseMode))
	if m.Payload != "" {
		buf.WriteString(fmt.Sprintf("%s[Payload]%s \n\n", tag, m.Payload))
	}
	if m.Args != nil {
		buf.WriteString(title + "[Args]")
		switch t := m.Args.(type) {
		case map[string]interface{}:
			for k, v := range t {
				buf.WriteString(fmt.Sprintf("%s %s: %v\n", tag, k, v))
			}
		case []interface{}:
			buf.WriteString(fmt.Sprintf("%s %s", tag, t))
		case string:
			buf.WriteString(tag + t)
		default:
			str, _ := jsoniter.MarshalToString(t)
			buf.WriteString(tag + str)
		}
	}
	if m.Error != nil {
		buf.WriteString(fmt.Sprintf("%s[Error]%s \n\n", tag, m.Error.Error()))
	}
	if m.Extra != nil {
		buf.WriteString(title + "[Extra]")
		for k, v := range m.Extra {
			buf.WriteString(fmt.Sprintf("%s %s: %v", tag, k, v))
		}
	}
	return buf.String()
}

func (m Message) String() string {
	b, _ := jsoniter.Marshal(m)
	return string(b)
}

func (c WebHookRebot) TxtReport(message ReportErrors) {
	var msg = Message{
		MsgType: TxtMessage,
		Text: struct {
			Content string `json:"content,omitempty"`
		}{
			Content: message.String(),
		},
	}
	_, err := http.Post(string(c), "application/json", strings.NewReader(msg.String()))
	if err != nil {
		logx.NewTraceLogger(context.Background()).Debug().Err(err).Msg("webhook POST失败")
	}
}

func (c WebHookRebot) MarkdownReport(message ReportErrors) {
	var msg = Message{
		MsgType: MarkdownMessage,
		Markdown: struct {
			Content string `json:"content,omitempty"`
		}{
			Content: message.MarkdownString(),
		},
	}
	_, err := http.Post(string(c), "application/json", strings.NewReader(msg.String()))
	if err != nil {
		logx.NewTraceLogger(context.Background()).Debug().Err(err).Msg("webhook POST失败")
	}
}
