package models

import (
	"context"
	"github.com/google/martian/log"
	"github.com/jmoiron/sqlx"
	"github.com/pengcainiao2/usercenter/internal/v1/form"
	"github.com/pengcainiao2/zero/tools/syncer"
	"time"
)

type CoinsAgentUsers struct {
	Uid        int64     `db:"uid" json:"uid"`
	Username   string    `db:"username" json:"username"`
	Password   string    `db:"password" json:"password"`
	Phone      string    `db:"phone" json:"phone"`
	Nickname   string    `db:"nickname" json:"nickname"`
	Gender     int64     `db:"gender" json:"gender"`
	Country    string    `db:"country" json:"country"` // ISO-3166 Country Code
	Salt       string    `db:"salt" json:"salt"`
	CreateAt   time.Time `db:"create_at" json:"create_at"`
	Banned     int64     `db:"banned" json:"banned"`
	Level      int64     `db:"level" json:"level"`
	UpLevelUid int64     `db:"up_level_uid" json:"up_level_uid"`
	UpNickname string    `db:"up_nickname" json:"up_nickname"`
	LimitNum   int64     `db:"limit_num" json:"limit_num"`
	DolaId     string    `db:"dola_id" json:"dola_id"`
	DelAt      int       `db:"del_at" json:"del_at"`
	DolaUid    string    `db:"dola_uid" json:"dola_uid"`
}

func GetUidRoles(ctx context.Context, uid string) (int64, error) {
	var level int64
	const authenticateUser = `SELECT level FROM okr.coins_agent_users WHERE uid = ? ;`
	err := syncer.MySQL().Get(ctx, &level, authenticateUser, uid)
	if err != nil {
		log.Errorf("GetUidRoles fail,err:=%v", err)
		return level, err
	}

	return level, err
}

func AuthenticateCoinAgentUser(ctx context.Context, in form.LoginReq) (uint64, error) {
	uid := 0
	const authenticateUser = `SELECT uid FROM okr.coins_agent_users WHERE username = ? AND password = MD5(CONCAT(?, salt));`
	err := syncer.MySQL().Get(ctx, &uid, authenticateUser, in.UserName, in.Password)
	return uint64(uid), err
}

func GetAgentUsersInfo(ctx context.Context, uids []uint64) (map[uint64]*CoinsAgentUsers, error) {
	var (
		list    []*CoinsAgentUsers
		dataMap = make(map[uint64]*CoinsAgentUsers, 0)
	)
	asql := `SELECT * FROM okr.coins_agent_users WHERE del_at = 0 AND uid in (?) ;`
	query, args, _ := sqlx.In(asql, uids)
	if err := syncer.MySQL().Select(ctx, &list, query, args...); err != nil {
		return dataMap, err
	}

	for _, v := range list {
		dataMap[uint64(v.Uid)] = v
	}
	return dataMap, nil
}
