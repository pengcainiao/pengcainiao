package syncer

import (
	"fmt"
	"strings"
	"sync"

	"gitlab.com/a16624741591/zero/core/env"
	"gitlab.com/a16624741591/zero/core/queue/nsqueue"
)

var (
	dataSyncerLock  = sync.Once{}
	dataSyncerQueue *nsqueue.Topic

	queueMap  = make(map[string]*nsqueue.Topic)
	queueLock sync.Mutex
)

func defaultPriority(defaultPriority nsqueue.Priority, priority ...nsqueue.Priority) nsqueue.Priority {
	if len(priority) == 0 {
		return defaultPriority
	}
	return priority[0]
}

func getQueueNsqTopic(name string, priority ...nsqueue.Priority) *nsqueue.Topic {
	name = fmt.Sprintf("%s%s", name, defaultPriority(nsqueue.NoOrderNormal, priority...))
	if !env.IsProduction() && !strings.HasSuffix(name, env.ReleaseMode) {
		name = fmt.Sprintf("%s-%s", name, env.ReleaseMode)
	}
	if queue, ok := queueMap[name]; ok {
		return queue
	}
	return nil
}

func QueueFlyele(priority ...nsqueue.Priority) *nsqueue.Topic {
	if topic := getQueueNsqTopic("flyele-nsq-flyelecore", priority...); topic != nil {
		return topic
	}

	//默认值 order_normal
	flyeleTopic := &nsqueue.Topic{
		Name:          "flyele-nsq-flyelecore",
		EnableOrdered: false,
		Channel:       "default_channel",
	}
	flyeleTopic.Priority(defaultPriority(nsqueue.NoOrderNormal, priority...))

	queueLock.Lock()
	queueMap[flyeleTopic.Name] = flyeleTopic
	queueLock.Unlock()

	return flyeleTopic
}

func QueueDataSyncer(priority ...nsqueue.Priority) *nsqueue.Topic {
	if topic := getQueueNsqTopic("flyele-nsq-datasyncer", priority...); topic != nil {
		return topic
	}

	dataSyncerTopic := &nsqueue.Topic{
		Name:          "flyele-nsq-datasyncer",
		EnableOrdered: false,
		Channel:       "default_channel",
	}
	dataSyncerTopic.Priority(defaultPriority(nsqueue.NoOrderNormal, priority...))

	queueLock.Lock()
	queueMap[dataSyncerTopic.Name] = dataSyncerTopic
	queueLock.Unlock()

	return dataSyncerTopic
}

func QueueCloudDisk(priority ...nsqueue.Priority) *nsqueue.Topic {
	if topic := getQueueNsqTopic("flyele-nsq-clouddisk", priority...); topic != nil {
		return topic
	}

	clouddiskTopic := &nsqueue.Topic{
		Name:          "flyele-nsq-clouddisk",
		EnableOrdered: false,
		Channel:       "default_channel",
	}
	clouddiskTopic.Priority(defaultPriority(nsqueue.NoOrderNormal, priority...))

	queueLock.Lock()
	queueMap[clouddiskTopic.Name] = clouddiskTopic
	queueLock.Unlock()

	return clouddiskTopic
}

func QueueUserCenter(priority ...nsqueue.Priority) *nsqueue.Topic {
	if topic := getQueueNsqTopic("flyele-nsq-usercenter", priority...); topic != nil {
		return topic
	}

	usercenterTopic := &nsqueue.Topic{
		Name:          "flyele-nsq-usercenter",
		EnableOrdered: false,
		Channel:       "default_channel",
	}
	usercenterTopic.Priority(defaultPriority(nsqueue.NoOrderNormal, priority...))

	queueLock.Lock()
	queueMap[usercenterTopic.Name] = usercenterTopic
	queueLock.Unlock()

	return usercenterTopic
}

func QueueUserInteraction(priority ...nsqueue.Priority) *nsqueue.Topic {
	if topic := getQueueNsqTopic("flyele-nsq-userinteraction", priority...); topic != nil {
		return topic
	}

	userinteractionTopic := &nsqueue.Topic{
		Name:          "flyele-nsq-userinteraction",
		EnableOrdered: false,
		Channel:       "default_channel",
	}
	userinteractionTopic.Priority(defaultPriority(nsqueue.NoOrderNormal, priority...))

	queueLock.Lock()
	queueMap[userinteractionTopic.Name] = userinteractionTopic
	queueLock.Unlock()

	return userinteractionTopic
}

func QueueRecord(priority ...nsqueue.Priority) *nsqueue.Topic {
	if topic := getQueueNsqTopic("flyele-nsq-record", priority...); topic != nil {
		return topic
	}

	recordTopic := &nsqueue.Topic{
		Name:          "flyele-nsq-record",
		EnableOrdered: false,
		Channel:       "default_channel",
	}
	recordTopic.Priority(defaultPriority(nsqueue.NoOrderNormal, priority...))

	queueLock.Lock()
	queueMap[recordTopic.Name] = recordTopic
	queueLock.Unlock()

	return recordTopic
}

func QueueTimedTask(priority ...nsqueue.Priority) *nsqueue.Topic {
	if topic := getQueueNsqTopic("flyele-nsq-timedtask", priority...); topic != nil {
		return topic
	}

	timedTaskTopic := &nsqueue.Topic{
		Name:          "flyele-nsq-timedtask",
		EnableOrdered: false,
		Channel:       "default_channel",
	}
	timedTaskTopic.Priority(defaultPriority(nsqueue.NoOrderNormal, priority...))

	queueLock.Lock()
	queueMap[timedTaskTopic.Name] = timedTaskTopic
	queueLock.Unlock()

	return timedTaskTopic
}

func QueuePushGatewayMulticast(priority ...nsqueue.Priority) *nsqueue.Topic {
	if topic := getQueueNsqTopic("flyele-nsq-pushgateway-multicast", priority...); topic != nil {
		return topic
	}

	pushGatewayMultTopic := &nsqueue.Topic{
		Name:          "flyele-nsq-pushgateway-multicast",
		Multicast:     true,
		EnableOrdered: false,
	}
	pushGatewayMultTopic.Priority(defaultPriority(nsqueue.NoOrderNormal, priority...))

	queueLock.Lock()
	queueMap[pushGatewayMultTopic.Name] = pushGatewayMultTopic
	queueLock.Unlock()

	return pushGatewayMultTopic
}

func QueuePushGatewayUnicast(priority ...nsqueue.Priority) *nsqueue.Topic {
	if topic := getQueueNsqTopic("flyele-nsq-pushgateway-unicast", priority...); topic != nil {
		return topic
	}

	pushGatewayUnicastTopic := &nsqueue.Topic{
		Name:          "flyele-nsq-pushgateway-unicast",
		EnableOrdered: false,
		Channel:       "default_channel",
	}
	pushGatewayUnicastTopic.Priority(defaultPriority(nsqueue.NoOrderNormal, priority...))

	queueLock.Lock()
	queueMap[pushGatewayUnicastTopic.Name] = pushGatewayUnicastTopic
	queueLock.Unlock()

	return pushGatewayUnicastTopic
}

func QueueStatistics(priority ...nsqueue.Priority) *nsqueue.Topic {
	if topic := getQueueNsqTopic("flyele-nsq-statistics", priority...); topic != nil {
		return topic
	}

	statisticsTopic := &nsqueue.Topic{
		Name:          "flyele-nsq-statistics",
		EnableOrdered: false,
		Channel:       "default_channel",
	}
	statisticsTopic.Priority(defaultPriority(nsqueue.NoOrderNormal, priority...))

	queueLock.Lock()
	queueMap[statisticsTopic.Name] = statisticsTopic
	queueLock.Unlock()

	return statisticsTopic
}

func InitialDataSyncerNoOrderNormal() *nsqueue.Topic {
	dataSyncerLock.Do(func() {
		dataSyncerQueue = QueueDataSyncer(nsqueue.NoOrderNormal)
		var (
			tps = nsqueue.Topics{dataSyncerQueue}
		)
		_ = tps.SetProducerManager()
	})
	return dataSyncerQueue
}
