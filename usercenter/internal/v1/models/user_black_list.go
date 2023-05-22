package models

import (
	"context"
	"github.com/pengcainiao/zero/core/discov"
	"github.com/pengcainiao/zero/core/env"
	"github.com/pengcainiao/zero/core/logx"
)

const blackList = "/pengcainiao/configs/user/black-users"

type UserBlackList struct{}

var (
	blackUserIDS []string
)

func init() {
	loadData()
}

func loadData() {
	etcd := EtcdCli()
	resp := etcd.LoadOrStore(blackList, "")
	_ = resp.JSON(&blackUserIDS)
	if len(blackUserIDS) == 0 {
		logx.NewTraceLogger(context.Background()).Warn().Msg("blackUserIDS is NULL")
	}
	go etcd.WatchKey(blackList, func(event discov.EtcdChangeType, value string) {
		logx.NewTraceLogger(context.Background()).Debug().Interface("event", event).Msg("watch event")
		_ = discov.GetResponse(value).JSON(&blackUserIDS)
	})
}

// IsBlackUser 用户是否是黑名单
func (u UserBlackList) IsBlackUser(ctx context.Context, userID string) bool {
	if !env.IsProduction() {
		return false
	}
	for _, uid := range blackUserIDS {
		if uid == userID {
			return true
		}
	}
	return false
}
