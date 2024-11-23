// NSQ 客户端
// 1、创建topic
//
//	curl -X POST 'http://nsq-lookup.infrastructure.svc.cluster.local:4161/topic/create?topic=flyele-nsq-socket-multicast-release&partition_num=2&replicator=1&extend=true&disable_channel_auto_create="true"&orderedmulti=true'
//
// 2、创建channel
// curl -X POST 'http://nsq-lookup.infrastructure.svc.cluster.local:4161/channel/create?topic=flyele-nsq-socket-multicast-release&channel=XXXX
package nsqueue

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"

	jsoniter "github.com/json-iterator/go"
	"github.com/spaolacci/murmur3"
	"github.com/youzan/go-nsq"
	"gitlab.com/a16624741591/zero/core/conf"
	"gitlab.com/a16624741591/zero/core/env"
	"gitlab.com/a16624741591/zero/core/logx"
	sonyflake "gitlab.com/a16624741591/zero/core/snowflake"
	"gitlab.com/a16624741591/zero/core/stores/redis"
	"gitlab.com/a16624741591/zero/core/timex"
	"gitlab.com/a16624741591/zero/core/trace"
	"gitlab.com/a16624741591/zero/rest/httprouter"
	"gitlab.com/a16624741591/zero/rest/httpx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	hook = httpx.WebHook()
)

// Topic
// 非顺序的topic, 由于支持同一个partition进行多个并发消费, 因此无需过多的partitions, 只需保证写入性能满足需求即可,
// 另外为了保持和原版nsq兼容, 每个节点只能有一个分区, 因此分区数*副本数不能大于节点总数. 非顺序的topic可以动态扩建分
// 区不影响业务使用. 建议普通topic使用 2分区2副本, 业务数据很多, 但是不怎么重要的, 比如log数据, 可以使用4分区1副本,
// 对于数据要求很高的, 可以使用2分区3副本(需要6台机器集群).
//
// 对于顺序的topic, 分区数取决于消费能力, 如果需要更高的并发度, 可以优化消费业务的消费能力, 如果消费能力已经无法优化,
// 则需要更多的分区数提高并发能力. 大部分情况下16个分区可以满足大部分需求, 如果吞吐量很大, 还需要实际计算消费业务的延
// 迟来决定. 由于顺序消费topic在动态扩建topic时会导致无序消息, 因此需要规划一个较长时间的潜在能力.
// 非顺序的topic，无需过多partition，能保证写入性能就行 ，建议2分区1副本
// 顺序的topic，建议2分区1副本
type Topic struct {
	Name                                  string //无需添加 env.ReleaseMode，会自动携带
	EnableOrdered                         bool   //是否顺序消费，优先于 EnablePartition
	Channel                               string //通道名称，producer可以不设置，让consumer去自动创建，默认值：default_channel
	MaxConsumeMQGoroutineCountPerConsumer int    // 单个consumer可以用来处理MQ消息的最大协程数，默认值：500
	Multicast                             bool   //是否为多播队列
	partitionNumber                       int    //有多少个分区，默认值：2
	replicator                            int    //有多少个副本，默认值：1
	syncDisk                              int    //多久刷磁盘，默认值：2000
	retentionDays                         int    //可以查看多久前的数据，默认值：3
	nsqmgr                                *Nsq
	initialLock                           sync.Once
	topicNameInitialed                    bool
}

type NsqDataProtocol struct {
	Topic     string        `json:"topic,omitempty"` //消息所属topic
	Body      interface{}   `json:"body,omitempty"`  //消息体
	TraceID   uint64        `json:"-"`               //traceID，仅最上游可不填写，其他层应传递上一级的traceID
	OrderedBy string        `json:"-"`               //仅当数据需要有序时提供，应提供如 task_id
	Header    MessageHeader `json:"header"`          //在传递中使用 json ext传输
}

// MessageHeader 消息头
type MessageHeader struct {
	Data     httprouter.HeaderData `json:"data"`
	CreateAt string                `json:"create_at"` //创建时间
}

func (h MessageHeader) toMap() map[string]interface{} {
	h.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	b, _ := json.Marshal(h)

	var d map[string]interface{}
	_ = json.Unmarshal(b, &d)
	return d
}

type nsqLogger struct {
}

func (l nsqLogger) Output(calldepth int, s string) error {
	if s == "" {
		return nil
	}
	fmt.Println("[NSQ]：", s)
	if strings.Contains(s, "No producer for") {
		go func() {
			hook.TxtReport(httpx.ReportErrors{
				Title:   "EMERGENCY",
				Payload: "NSQ topic出现问题，将自动重启topic",
				Args: map[string]interface{}{
					"message": s,
				},
			})
		}()
	} else if !strings.HasPrefix(s, "INF") {
		hook.TxtReport(httpx.ReportErrors{
			Title:   "EMERGENCY",
			Payload: "NSQ 出现问题，请检查队列情况，并马上通知维护人员",
			Args: map[string]interface{}{
				"message": s,
			},
		})
	}
	return nil
}

type Nsq struct {
	Endpoint                 string
	maxConsumeGoroutineCount int
	config                   *nsq.Config           //配置信息
	topicMgr                 *nsq.TopicProducerMgr //topic
	alreadyRegistered        map[string][]string   //
}

func (np NsqDataProtocol) Bytes() []byte {
	b, _ := json.Marshal(np)
	return b
}

func (np *NsqDataProtocol) extraTrace(producer bool, topic *Topic) (*httprouter.Context, oteltrace.Span) {
	if np.Header.Data.RequestID != "" {
		ctx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.MapCarrier{trace.OlteTraceHeader: np.Header.Data.RequestID})
		tracer := trace.GetTracerProvider(trace.NsqTraceName)
		if producer {
			_, span := tracer.Start(ctx, "[NSQ Producer]"+topic.Name,
				oteltrace.WithSpanKind(oteltrace.SpanKindProducer),
			)
			span.SetAttributes(
				attribute.String(string(semconv.MessagingSystemKey), "nsq"),
				attribute.Int64(string(semconv.MessageIDKey), int64(np.TraceID)))
			trace.PatchSpanData(span)
			return nil, span
		} else {
			spanContext, span := tracer.Start(ctx, "[NSQ Consumer] "+topic.Name,
				oteltrace.WithSpanKind(oteltrace.SpanKindConsumer),
			)
			trace.PatchSpanData(span)
			d, _ := json.MarshalToString(np.Body)
			span.AddEvent("message", oteltrace.WithAttributes(
				attribute.String("message.data", d),
				attribute.String(string(semconv.MessagingSystemKey), "nsq"),
				attribute.Int64(string(semconv.MessageIDKey), int64(np.TraceID)),
			))
			var mapCarrier = propagation.MapCarrier{}
			propagator := otel.GetTextMapPropagator()
			propagator.Inject(spanContext, mapCarrier)
			np.Header.Data.RequestID = mapCarrier[trace.OlteTraceHeader]
			return httprouter.NewContextData(oteltrace.ContextWithSpan(ctx, span), &httprouter.HeaderData{
				RequestID: np.Header.Data.RequestID,
			}), span
		}
	}
	return httprouter.NewContextData(context.Background(), nil), nil
}

func FromNsqMessage(message *nsq.Message, consume MqMessageConsumer) *NsqDataProtocol {
	var header = MessageHeader{}
	g, err := message.GetJsonExt()
	if err != nil {
		logx.NewTraceLogger(context.Background()).Err(err).Msg("尝试从消息中获取ext时失败")
	} else {
		b, _ := json.Marshal(g.Custom)
		_ = json.Unmarshal(b, &header)
	}
	np := UnmarshalNsqMessage(message, &consume)
	//np.Body = consume
	np.Header = header
	return np
}

// UnmarshalNsqMessage 将NSQ的消息转换为指定的对象
func UnmarshalNsqMessage(msg *nsq.Message, bodyDst *MqMessageConsumer) *NsqDataProtocol {
	var np = &NsqDataProtocol{Body: bodyDst}
	_ = json.Unmarshal(msg.Body, &np)

	// 兼容赋值问题
	var newNp NsqDataProtocol
	_ = json.Unmarshal(msg.Body, &newNp)
	logx.NewTraceLogger(context.Background()).Debug().Interface("body new np", newNp).Msg("")
	return &newNp
}

func (t *Topic) check() {
	if !env.IsProduction() && !strings.HasSuffix(t.Name, env.ReleaseMode) {
		t.Name = fmt.Sprintf("%s-%s", t.Name, env.ReleaseMode)
	}
	t.partitionNumber = 2
	if t.replicator > 1 || t.replicator == 0 {
		t.replicator = 1
	}
	if t.syncDisk < 2000 {
		t.syncDisk = 2000
	}
	if t.retentionDays <= 0 || t.retentionDays > 5 {
		t.retentionDays = 3
	}
}

// NewNsqLookup 创建nsqd实例
func NewNsqLookup(opts ...conf.Option) *Nsq {
	cfg := nsq.NewConfig()
	cfg.Hasher = murmur3.New32()
	cfg.EnableTrace = !env.IsProduction()
	cfg.LowRdyTimeout = time.Second * 5
	cfg.EnableTrace = true
	var n = &Nsq{
		config:            cfg,
		alreadyRegistered: map[string][]string{},
	}
	var lcfg = conf.ApplyConfig(opts...)
	n.Endpoint = lcfg.NsqAddress
	n.config.MaxInFlight = lcfg.MaxInFlight
	n.maxConsumeGoroutineCount = lcfg.MaxNSQGoroutineCount
	return n
}

func (ns *Nsq) ensureTopic(topics ...*Topic) []string {
	var t = make([]string, 0)
	for _, topic := range topics {
		topic.check()
		t = append(t, topic.Name)
	}
	return t
}

// NewProducer 创建消费者
func (ns *Nsq) NewProducer(topic ...*Topic) error {
	if !env.IsRunningInK8s() {
		return errNotRunningWithK8s
	}

	var topics = ns.ensureTopic(topic...)
	rand.Seed(time.Now().UnixNano())
	pubMgr, err := nsq.NewTopicProducerMgr(topics, ns.config)
	if err != nil {
		logx.NewTraceLogger(context.Background()).Err(err).Strs("topics", topics).Msg("NewTopicProducerMgr：error")
		return err
	}
	pubMgr.SetLogger(&nsqLogger{}, nsq.LogLevelInfo)
	var lookupAddress = ns.Endpoint
	if err := pubMgr.ConnectToNSQLookupd(lookupAddress); err != nil {
		logx.NewTraceLogger(context.Background()).Err(err).Str("connect_to", lookupAddress).Msg("ConnectToNSQLookupd：error")
		return err
	}
	ns.topicMgr = pubMgr
	return nil
}

// Publish 发布消息
func (ns *Nsq) Publish(topic *Topic, protocol NsqDataProtocol) error {
	if !env.IsRunningInK8s() {
		return errNotRunningWithK8s
	}
	if ns.topicMgr == nil {
		return errors.New("topicmgr is nil,invoke NewProducer first")
	}
	var (
		err    error
		msgExt = &nsq.MsgExt{
			TraceID: protocol.TraceID,
			Custom:  protocol.Header.toMap(),
		}
	)
	if protocol.Header.Data.RequestID == "" {
		logx.NewTraceLogger(context.Background()).Warn().Str("topic", topic.Name).Msg("nsq 未设置x-request-id")
	}
	if msgExt.TraceID == 0 {
		msgExt.TraceID = sonyflake.GenerateInt64ID()
		protocol.TraceID = msgExt.TraceID
	}
	_, span := protocol.extraTrace(true, topic)
	if span != nil {
		defer span.End()
	}

	if topic.EnableOrdered {
		_, _, _, err = ns.topicMgr.PublishOrderedWithJsonExt(topic.Name, []byte(protocol.OrderedBy), protocol.Bytes(), msgExt)
	} else {
		_, _, _, err = ns.topicMgr.PublishWithJsonExtAndPartitionId(topic.Name, rand.Intn(9)%int(topic.partitionNumber),
			protocol.Bytes(), msgExt)
	}
	return err
}

// NewConsumer 创建新的消息队列消费者
func (ns *Nsq) NewConsumer(topic *Topic, handler nsq.Handler) error {
	if topic.Channel == "" {
		topic.Channel = "default_channel"
	}
	ns.ensureTopic(topic)
	ns.config.EnableOrdered = topic.EnableOrdered
	ns.config.DefaultRequeueDelay = 10 * time.Second // 10秒重试一次

	if handler == nil {
		return errors.New("handler is nil")
	}
	if topic.EnableOrdered {
		cons, err := nsq.NewConsumer(topic.Name, topic.Channel, ns.config)
		if err != nil {
			logx.NewTraceLogger(context.Background()).Err(err).Msg("NewConsumer 失败")
			return err
		}
		if err := ns.startConsume(topic, handler, cons); err != nil {
			return err
		}
	} else {
		for paritonID := 0; paritonID < topic.partitionNumber; paritonID++ {
			cons, err := nsq.NewPartitionConsumer(topic.Name, paritonID, topic.Channel, ns.config)
			if err != nil {
				logx.NewTraceLogger(context.Background()).Err(err).Msg("NewPartitionConsumer 失败")
				return err
			}
			if err = ns.startConsume(topic, handler, cons); err != nil {
				return err
			}
		}
	}
	select {}
}

func (ns *Nsq) startConsume(topic *Topic, handler nsq.Handler, cons *nsq.Consumer) error {
	var (
		nsqConns     = cons.Stats()
		clientCounts int
	)
	if nsqConns == nil || nsqConns.Connections == 0 {
		clientCounts = 1
	}

	cons.ChangeMaxInFlight(clientCounts * topic.partitionNumber * ns.config.MaxInFlight)
	cons.SetConsumeExt(true)
	cons.AddConcurrentHandlers(handler, ns.maxConsumeGoroutineCount)

	if err := cons.ConnectToNSQLookupd(ns.Endpoint); err != nil {
		logx.NewTraceLogger(context.Background()).Err(err).Msg("ConnectToNSQLookupd 失败")
		return err
	}
	return nil
}

// IsUniqueMessageInSeconds 验证消息ID在一段时间内是否唯一，当前逻辑存在消息消费不准确的情况，使用 IsUniqueMessage 代替，3小时内确保唯一
// Deprecated
func IsUniqueMessageInSeconds(traceID uint64, seconds int) bool {
	if traceID < 10000000000 {
		// 日程提醒的traceID生成算法：
		// crc32(creator_id)+ time.Now().Format("20060102")
		return true
	}
	if sonyflake.ElapsedSeconds(traceID) > int64(seconds) {
		return false
	}
	var key = fmt.Sprintf("extraTrace:unq:%d", traceID)
	return redis.Client().SetNX(context.Background(), key, 1, time.Second*time.Duration(seconds)).Val()
}

// IsUniqueRdsKey 识别消息是否唯一
func IsUniqueRdsKey(rdsKey string) (unique bool) {
	return redis.Client().SetNX(context.Background(), rdsKey, "1", timex.RandomExpireSeconds(time.Minute*5)).Val()
}
