package syncer

import (
	"context"
	"errors"

	json "github.com/json-iterator/go"
	"gitlab.com/a16624741591/zero/core/logx"
	"gitlab.com/a16624741591/zero/core/stores/redis"
)

const (
	//WebsocketChannel ExchangeChannel = "flyele-websocket-exchange"
	//SseChannel       ExchangeChannel = "flyele-sse-exchange"

	Normal ChannelType = 2
	SSE    ChannelType = 1
)

type (
	ExchangeChannel string
	ChannelType     int8
)

type PushChannel struct {
	Channel     string      `json:"channel,omitempty"`
	ChannelType ChannelType `json:"channel_type,omitempty"`
}

type ExchangeMessage struct {
	PushChannel      *PushChannel           `json:"push_channel,omitempty"` //推送通道
	Code             int                    `json:"code,omitempty"`
	NotifyTo         []string               `json:"notify_to,omitempty"`
	Message          string                 `json:"message,omitempty"`
	Title            string                 `json:"title,omitempty"`
	RefID            string                 `json:"ref_id,omitempty"`
	RefType          string                 `json:"ref_type,omitempty"` // 消息类型
	FileID           string                 `json:"file_id,omitempty"`
	Changes          map[string]interface{} `json:"changes,omitempty"`    // 变更信息
	IsRedDot         bool                   `json:"is_red_dot,omitempty"` // 是否红点
	CommentID        string                 `json:"comment_id,omitempty"`
	CreatorID        string                 `json:"creator_id,omitempty" `  //消息创建人ID
	AssociateID      string                 `json:"associate_id,omitempty"` //消息关联ID
	SendFromPlatform string                 `json:"send_from,omitempty"`    //消息发送方平台
	AffectedUID      []string               `json:"affected_uid"`           // 受影响用户
}

func (ex ExchangeMessage) String() string {
	b, _ := json.Marshal(ex)
	return string(b)
}

func (ex ExchangeChannel) Publish(message ExchangeMessage, sse bool) error {
	if sse {
		if message.PushChannel == nil {
			return errors.New("push_channel is required")
		}
		message.PushChannel.ChannelType = SSE
	}
	err := redis.RedisClient().Publish(context.Background(), string(ex), message.String()).Err()
	if err != nil {
		logx.NewTraceLogger(context.Background()).Err(err).Str("channel", string(ex)).Msg("【pubsub】发布消息失败")
	}
	return err
}

func (ex ExchangeChannel) ConsumeMessage(onMessage func(channel, message string)) {
	channel := string(ex)
	pubsub := redis.RedisClient().Subscribe(context.Background(), channel)
	defer func() {
		_ = pubsub.Close()
	}()

	_, err := pubsub.Receive(context.Background())
	if err != nil {
		logx.NewTraceLogger(context.Background()).Err(err).Str("channel", string(ex)).Msg("【pubsub】接收消息失败")
		return
	}

	var pubch = pubsub.Channel()
	for message := range pubch {
		onMessage(channel, message.Payload)
	}
}
