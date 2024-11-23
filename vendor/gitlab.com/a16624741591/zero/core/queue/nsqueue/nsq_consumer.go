package nsqueue

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"gitlab.com/a16624741591/zero/core/sysx"

	"gitlab.com/a16624741591/zero/core/env"

	"gitlab.com/a16624741591/zero/core/stores/redis"
	"gitlab.com/a16624741591/zero/core/stores/sqlx"

	"gitlab.com/a16624741591/zero/rest/httpx"

	"github.com/youzan/go-nsq"
	"gitlab.com/a16624741591/zero/core/logx"
	"gitlab.com/a16624741591/zero/rest/httprouter"
)

const NsqTraceID string = "nsq_trace_id"

var (
	hostname = os.Getenv("HOSTNAME")
)

type MqMessageConsumer interface {
	HandleMessage(context *httprouter.Context, protocol *NsqDataProtocol) error
	New() MqMessageConsumer
}

type wrappedMqMessageConsumerHandler struct {
	//antsPool *ants.Pool        //协程池
	Topic    *Topic            //topic信息
	Consumer MqMessageConsumer //consumer
}

func (w wrappedMqMessageConsumerHandler) HandleMessage(message *nsq.Message) error {
	//if w.antsPool == nil {
	//w.antsPool, _ = ants.NewPool(500)
	//}
	//_ = w.antsPool.Submit(func() {
	var (
		err     error
		traceID = message.GetTraceID()
		rdsKey  = w.getUniqueMessageID(traceID)
	)

	defer func() {
		if err != nil {
			w.consumeFailed(message, err)
			logx.NewTraceLogger(context.Background()).Err(err).Uint64("trace_id", traceID).Msg("消息处理失败")
		} else {
			message.Finish()
		}
	}()

	if !IsUniqueRdsKey(rdsKey) {
		logx.NewTraceLogger(context.Background()).Error().Uint64("nsq_trace_id", traceID).
			Str("data", string(message.Body)).
			Str("topic", w.Topic.Name).
			Msg("消息已经处理过了")
		return nil
	}

	if w.Consumer != nil {
		consumer := w.Consumer.New()
		np := FromNsqMessage(message, consumer)
		np.TraceID = traceID
		ctx, span := np.extraTrace(false, w.Topic)
		if span != nil {
			defer span.End()
		}
		var httpctx = &httprouter.Context{
			Context: ctx,
			Data:    np.Header.Data,
		}
		err = consumer.HandleMessage(httpctx, np)
	} else {
		err = errors.New("未设置消息处理程序 w.Consumer == nil")
	}
	//})

	return err
}

func (w wrappedMqMessageConsumerHandler) getUniqueMessageID(traceID uint64) string {
	if w.Topic != nil && w.Topic.Multicast {
		if hostname == "" {
			hostname = os.Getenv("HOSTNAME")
		}
		return "extraTrace:unq:" + hostname + strconv.FormatUint(traceID, 10)
	}
	return "extraTrace:unq:" + sysx.SubSystem + strconv.FormatUint(traceID, 10)
}

func (w wrappedMqMessageConsumerHandler) consumeFailed(message *nsq.Message, err error) {
	np := FromNsqMessage(message, w.Consumer)
	np.TraceID = message.GetTraceID()

	failRedisKey := fmt.Sprintf("nsq:fail:%d", np.TraceID)
	failCount := redis.Client().Incr(context.Background(), failRedisKey).Val()
	if failCount == 1 {
		_ = redis.Client().Expire(context.Background(), failRedisKey, 2*time.Minute).Err()
	}

	if failCount == 3 {
		message.Finish()

		_ = redis.Client().Del(context.Background(), failRedisKey).Err()

		var sql = `INSERT INTO fx_mq_retry(trace_id,topic,message_create_at, data,error) VALUES(?,?,?,?,?)`
		_, _ = sqlx.NewMysql(env.DbDSN).Exec(sql, np.TraceID, w.Topic.Name, np.Header.CreateAt, string(np.Bytes()), err.Error())

		httpx.WebHook().MarkdownReport(httpx.ReportErrors{
			Title: "有一条来自消息队列的消息无法处理，已存储到数据库，需人工介入!",
			Args: map[string]interface{}{
				"data":        string(message.Body),
				"topic":       w.Topic.Name,
				"err_msg":     err.Error(),
				"trace_id":    np.TraceID,
				"server_name": sysx.SubSystem,
			},
		})
	}
}
