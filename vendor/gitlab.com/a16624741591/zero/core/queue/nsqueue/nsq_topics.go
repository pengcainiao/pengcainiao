package nsqueue

import (
	"context"
	"errors"

	"github.com/youzan/go-nsq"
	"gitlab.com/a16624741591/zero/core/conf"
	"gitlab.com/a16624741591/zero/core/env"
	"gitlab.com/a16624741591/zero/core/logx"
	"gitlab.com/a16624741591/zero/rest/httprouter"
)

var (
	errNotRunningWithK8s = errors.New("服务未运行在k8s环境中")
	//NsqDataSyncerSvcTopic 消费elasticsearch，MySQL数据
)

type (
	Topics   []*Topic //NSQ队列主体
	Priority string   //优先级
)

const (
	NoOrderNormal Priority = "-norder"
	NoOrderLow    Priority = "-norder-low"
	NoOrderHigh   Priority = "-norder-high"
	OrderedNormal Priority = ""     //正常优先级，使用不同的字符串进行标记并进行统一处理
	OrderedLow    Priority = "-low" //低优先级，一般限制处理速率
	OrderedHigh   Priority = "-high"
)

// SetProducerManager 一个服务中消费多个topic时，先创建NSQ实例，减少连接数
//
//	example
//	var topic = &Topic{
//		Name:            "nsq-test-topic",
//		EnableOrdered:   false,
//		EnablePartition: true,
//		Channel:         "default_channel",
//	}
//	// 发送数据
//	topic.Publish(NsqDataProtocol{
//		Topic:   topic.Name,
//		Body:    map[string]interface{}{"abc": 123, "b": "454545"},
//		RequestID: tools.GenerateInt64ID(),
//	})
//	//消费数据
//	topic.StartConsume(&testConsumerHandler{})
func (t Topics) SetProducerManager(opts ...conf.Option) error {
	if !env.IsRunningInK8s() {
		return errNotRunningWithK8s
	}
	var msgr = NewNsqLookup(opts...)
	if err := msgr.NewProducer(t...); err != nil {
		return err
	}
	for _, tp := range t {
		tp.nsqmgr = msgr
	}
	return nil
}

func (t *Topic) Priority(priority Priority) *Topic {
	if t.topicNameInitialed {
		return t
	}
	if priority == OrderedLow || priority == OrderedNormal || priority == OrderedHigh {
		t.EnableOrdered = true
	} else {
		t.EnableOrdered = false
	}
	t.Name += string(priority)
	t.check()
	t.topicNameInitialed = true
	return t
}

// Producer 向topic发送数据，设置header中的 SentryTraceID
func (t *Topic) Producer(ctx *httprouter.Context, data NsqDataProtocol, opts ...conf.Option) error {
	if !env.IsRunningInK8s() {
		logx.NewTraceLogger(context.Background()).Err(errNotRunningWithK8s).Msg("publish 消息失败")
		return nil
	}
	var err error
	if t.nsqmgr == nil {
		t.initialLock.Do(func() {
			t.nsqmgr = NewNsqLookup(opts...)
			err = t.nsqmgr.NewProducer(t)
		})
	}
	if err != nil {
		return err
	}

	if ctx != nil {
		data.Header.Data = ctx.Data
	}
	return t.nsqmgr.Publish(t, data)
}

// Consume 使用内部包装的类型 wrappedMqMessageConsumerHandler 统一处理了使用MQ要面临的分布式问题，包括幂等、失败后的重试策略等
// 参数 consumer 可不使用指针，但 New()方法返回值必须为指针
func (t *Topic) Consume(consumer MqMessageConsumer, options ...conf.Option) error {
	var err error
	if t.nsqmgr == nil {
		t.initialLock.Do(func() {
			t.nsqmgr = NewNsqLookup(options...)
			err = t.nsqmgr.NewProducer(t)
		})
	}
	if err != nil {
		return err
	}
	//if t.MaxConsumeMQGoroutineCountPerConsumer == 0 {
	//t.MaxConsumeMQGoroutineCountPerConsumer = 500
	//}
	//gopool, _ := ants.NewPool(t.MaxConsumeMQGoroutineCountPerConsumer)
	var wrapped = wrappedMqMessageConsumerHandler{
		//antsPool: gopool,
		Consumer: consumer,
		Topic:    t,
	}
	return t.nsqmgr.NewConsumer(t, wrapped)
}

// StartConsume 开始消费数据
// Deprecated: 使用 Consume 方法代替，新的方法封装了复杂过程
func (t *Topic) StartConsume(handler nsq.Handler, opts ...conf.Option) error {
	if !env.IsRunningInK8s() {
		return errNotRunningWithK8s
	}
	var err error
	if t.nsqmgr == nil {
		t.initialLock.Do(func() {
			t.nsqmgr = NewNsqLookup(opts...)
			err = t.nsqmgr.NewProducer(t)
		})
	}
	if err != nil {
		return err
	}
	return t.nsqmgr.NewConsumer(t, handler)
}

// Publish 向topic发送数据，设置header中的 SentryTraceID
// Deprecated:使用 Producer 方法代替
func (t *Topic) Publish(data NsqDataProtocol, opts ...conf.Option) error {
	if !env.IsRunningInK8s() {
		logx.NewTraceLogger(context.Background()).Err(errNotRunningWithK8s).Msg("publish 消息失败")
		return nil
	}
	var err error
	if t.nsqmgr == nil {
		t.initialLock.Do(func() {
			t.nsqmgr = NewNsqLookup(opts...)
			err = t.nsqmgr.NewProducer(t)
		})
	}
	if err != nil {
		return err
	}
	return t.nsqmgr.Publish(t, data)
}
