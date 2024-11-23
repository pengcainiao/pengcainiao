package discov

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"gitlab.com/a16624741591/zero/core/logx"

	"gitlab.com/a16624741591/zero/core/conf"

	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	PUT    EtcdChangeType = "put"
	DELETE EtcdChangeType = "delete"
)

var (
	defaultTimeout = time.Second * 2
	etcdClient     *EtcdClient
)

type EtcdChangeType string
type GetResponse string

// WatchKeyChanged 当监听的key发生变化时发送
type WatchKeyChanged func(event EtcdChangeType, v string)

type EtcdClient struct {
	*clientv3.Client
}

func Etcd(opts ...conf.Option) *EtcdClient {
	if etcdClient == nil {
		var s = newEtcdClient(opts...)
		etcdClient = &s
	}
	return etcdClient
}

func newEtcdClient(opts ...conf.Option) EtcdClient {
	cfg := conf.ApplyConfig(opts...)
	ec, err := clientv3.New(clientv3.Config{
		Endpoints: []string{cfg.EtcdAddress},
	})
	if err != nil {
		return EtcdClient{}
	}
	return EtcdClient{
		ec,
	}
}

func (t *EtcdClient) GetRegisteredKey(keyName string) (map[string]string, error) {
	rangeResp, err := t.Get(context.TODO(), keyName, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var m = make(map[string]string)
	for _, kv := range rangeResp.Kvs {
		//t.data[string(kv.Key)] = string(kv.Value)
		m[string(kv.Key)] = string(kv.Value)
	}

	go t.WatchKey(keyName, nil)
	return m, nil
}

// WatchKey 监控服务目录下的事件
func (t *EtcdClient) WatchKey(keyName string, onChange WatchKeyChanged) {
	// Watch 服务目录下的更新
	watchChan := t.Watch(context.TODO(), keyName, clientv3.WithPrefix())
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			if onChange != nil {
				var t = PUT
				if event.Type == mvccpb.DELETE {
					t = DELETE
				}
				onChange(t, string(event.Kv.Value))
			}
		}
	}
}

// LoadOrStore 加载或存储
func (t *EtcdClient) LoadOrStore(keyName, value string, opts ...clientv3.OpOption) GetResponse {
	ctx, cancel := context.WithTimeout(context.TODO(), defaultTimeout)
	go func() {
		time.Sleep(defaultTimeout)
		cancel()
	}()
	gresp, err := t.Get(ctx, keyName, opts...)
	if err != nil {
		logx.NewTraceLogger(context.Background()).Err(err).Str("key", keyName).Msg("读取ETCD失败")
		return ""
	}
	if gresp != nil && len(gresp.Kvs) > 0 {
		for _, kv := range gresp.Kvs {
			return GetResponse(kv.Value)
		}
	} else {
		_, _ = t.Put(context.Background(), keyName, value, opts...)
		return ""
	}
	return ""
}

func (t *EtcdClient) Register(keyName string, callback func()) {
	var (
		curLeaseId   clientv3.LeaseID = 0
		callbackOnce                  = sync.Once{}
	)
	for {
		if curLeaseId == 0 {
			leaseResp, err := t.Grant(context.TODO(), 5)
			if err != nil {
				logx.NewTraceLogger(context.Background()).Err(err).Msg("【etcd】grant失败")
				continue
			}

			key := keyName + fmt.Sprintf("-%d", leaseResp.ID)
			if _, err := t.Put(context.TODO(), key, "1", clientv3.WithLease(leaseResp.ID)); err != nil {
				logx.NewTraceLogger(context.Background()).Err(err).Str("put_key", key).Msg("【etcd】创建key失败")
				continue
			}
			callbackOnce.Do(func() {
				callback()
			})
			curLeaseId = leaseResp.ID
		} else {
			// 续约租约，如果租约已经过期将curLeaseId复位到0重新走创建租约的逻辑
			if _, err := t.KeepAliveOnce(context.TODO(), curLeaseId); err == rpctypes.ErrLeaseNotFound {
				curLeaseId = 0
				continue
			}

		}
		time.Sleep(time.Second)
	}
}

func (g GetResponse) JSON(dest interface{}) error {
	var t = reflect.TypeOf(dest)
	if t.Kind() != reflect.Ptr {
		return errors.New("dest should be ptr")
	}
	return json.Unmarshal([]byte(g), &dest)
}

func (g GetResponse) String() string {
	return string(g)
}

func (g GetResponse) Int64() (int64, error) {
	return strconv.ParseInt(string(g), 10, 64)
}
