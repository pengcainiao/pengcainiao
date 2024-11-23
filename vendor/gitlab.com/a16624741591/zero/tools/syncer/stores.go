package syncer

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pengcainiao/sqlx"
	"gitlab.com/a16624741591/zero/core/discov"
	"gitlab.com/a16624741591/zero/core/env"
	"gitlab.com/a16624741591/zero/core/logx"
	"gitlab.com/a16624741591/zero/core/queue/nsqueue"
	sonyflake "gitlab.com/a16624741591/zero/core/snowflake"
	"gitlab.com/a16624741591/zero/core/stores/redis"
	coresqlx "gitlab.com/a16624741591/zero/core/stores/sqlx"
	_ "gitlab.com/a16624741591/zero/core/sysx"
	"gitlab.com/a16624741591/zero/rest/httpx"
)

type HoldingType string

const (
	HoldingUsers HoldingType = "bloom:user" //HoldingUsers 存储用户ID
	//HoldingTasks                HoldingType = "bloom:task"      //HoldingTasks 存储事项或会议ID
	//HoldingTaskDispatches       HoldingType = "bloom:task:disp" //HoldingTaskDispatches 存储分发ID
	HoldingRecords              HoldingType = "bloom:record"  //HoldingRecords 存储记录ID
	HoldingProjects             HoldingType = "bloom:project" //HoldingProjects 存储评论ID
	HoldingFiles                HoldingType = "bloom:files"   //HoldingFiles 存储文件ID
	HoldingComments             HoldingType = "bloom:commemt" //HoldingComments 存储评论ID
	DefaultExpireDuration                   = time.Hour * 24  // DefaultExpireDuration 默认过期时间
	NotExistsItemExpireDuration             = time.Hour * 2   // NotExistsItemExpireDuration 查询不到过期时间
)

func init() {
	sqlx.SetIncrementalFoundedCallback(func(payload string, args ...interface{}) {
		var m map[string]interface{}
		_ = json.Unmarshal([]byte(payload), &m)

		var err = InitialDataSyncerNoOrderNormal().Publish(nsqueue.NsqDataProtocol{
			Topic:   "mysql-incremental",
			Body:    m,
			TraceID: sonyflake.GenerateInt64ID(),
		})
		if err != nil {
			_, _ = MySQL().Exec(context.Background(), "INSERT INTO incre_errors(full_data,reason) VALUES(?,?)", payload, err.Error())

			var reportErr = httpx.ReportErrors{
				Payload: payload,
				Args:    args,
				Error:   err,
			}

			httpx.WebHook().MarkdownReport(reportErr)
		}

	})
}

// MySQL mysql数据库连接
func MySQL() *coresqlx.SqlxDB {
	return coresqlx.MySQL()
}

// Redis redis数据库连接
func Redis() redis.RedisNode {
	return redis.Client()
}

// Etcd ETCD客户端
func Etcd() *discov.EtcdClient {
	return discov.Etcd()
}

// oss oss客户端
func Oss() *oss.Client {
	aliKeyID := env.AliResourcesAccessKey
	aliKeySecret := env.AliResourcesAccessSecret
	client, err := oss.New("http://oss-cn-shenzhen.aliyuncs.com", aliKeyID, aliKeySecret)
	if err != nil {
		logx.NewTraceLogger(context.Background()).Err(err).Msg("初始化oss客户端失败")
	}
	return client
}

// DefaultExpireTime 默认过期时间
func DefaultExpireTime(isNotFound ...bool) time.Duration {
	if len(isNotFound) > 0 && isNotFound[0] {
		return RandomExpireSeconds(NotExistsItemExpireDuration)
	}
	return RandomExpireSeconds(time.Duration(26-time.Now().Hour()) * time.Hour)
}

// DefaultExpireTimeHour 默认过期时间小时
func DefaultExpireTimeHour(hour ...int) time.Duration {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(600)
	if len(hour) > 0 {
		return time.Duration(n) + time.Duration(hour[0])*time.Hour
	}
	return time.Duration(n) + time.Hour
}

// DefaultExpireTimeMinute 默认过期时间分钟
func DefaultExpireTimeMinute(minute ...int) time.Duration {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10)
	if len(minute) > 0 {
		return time.Duration(n) + time.Duration(minute[0])*time.Minute
	}
	return time.Duration(n) + time.Minute
}

// RandomExpireSeconds 在基准值的基础上加2小时随机时间
func RandomExpireSeconds(base time.Duration) time.Duration {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(1800)
	return time.Second*time.Duration(n) + base
}
