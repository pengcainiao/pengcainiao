package models

import (
	"github.com/pengcainiao/pengcainiao/usercenter/internal/v1/form"
	"github.com/pengcainiao/zero/rest/httprouter"
)

//UserOnline 在线用户信息
type UserOnline struct {
	form.DeviceInfoV2
	WechatSessionKey string `json:"-"`
	Account          string `json:"-" db:"account"`
	UnionID          string `json:"-"`
	UserID           string `json:"user_id" db:"user_id"`                 //用户ID
	LastLoginTime    int64  `json:"last_login_time" db:"last_login_time"` //最后一次登录时间
	AppID            string //用户登录进来时使用的小程序ID
	Online           int
	//DeviceID         string `json:"device_id" db:"device_id"`
}

//UserOnline 用户上线
func (u UserOnline) UserOnline(ctx *httprouter.Context) string {
	var token = NewTokenObject()
	token.ClientVersion = u.ClientVersion
	_ = token.New(ctx, u, "")
	return token.Token
}
