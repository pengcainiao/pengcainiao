package models

import (
	"github.com/pengcainiao/pengcainiao/usercenter/internal/v1/form"
	"github.com/pengcainiao/zero/core/logx"
	sonyflake "github.com/pengcainiao/zero/core/snowflake"
	"github.com/pengcainiao/zero/rest/httprouter"
	"github.com/pengcainiao/zero/tools/syncer"
	"time"
)

type User struct {
	ID       string `json:"id" db:"id"`
	UserName string `json:"user_name" db:"user_name"`
	Password string `json:"password" db:"password"`
}

func (u *User) UserIsExist(ctx *httprouter.Context, userName string) bool {
	var user []User
	asql := "SELECT ID FROM user WHERE user_name = ? AND status =1;"
	err := syncer.MySQL().Select(ctx, &user, asql, userName)
	if err != nil {
		logx.NewTraceLogger(ctx).Err(err).Interface("userName:", userName).Msg("UserIsExist sql fail")
		return false
	}
	if len(user) != 0 {
		u.ID = user[0].ID
		return true
	} else {
		return false
	}
}

func (u User) RegisterUser(ctx *httprouter.Context, params form.RegisterUserRequest) (form.RegisterAndLoginResponse, error) {
	var data form.RegisterAndLoginResponse
	asql := "INSERT INTO user (id, user_name, password) VALUES (?,?,?)"
	id := sonyflake.GenerateID()
	_, err := syncer.MySQL().Exec(ctx, asql, id, params.UserName, params.Password)
	if err != nil {
		logx.NewTraceLogger(ctx).Err(err).Interface("params", params).Msg("RegisterUser sql fail")
		return form.RegisterAndLoginResponse{}, err
	}
	BloomAddNewUser(id)
	data.UserID = id

	return data, nil
}

func (u *User) UserLogin(ctx *httprouter.Context, params form.LoginRequest) string {
	onlineUser := UserOnline{
		UserID:        u.ID,
		LastLoginTime: time.Now().Unix(),
		DeviceInfoV2: form.DeviceInfoV2{
			Platform:      "pc",
			ClientVersion: "16",
			DeviceID:      "dc5c4057a68e6c880107fddb237f0775",
		},
	}
	token := onlineUser.UserOnline(ctx)
	return token
}
