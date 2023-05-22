package models

import (
	"github.com/pengcainiao/pengcainiao/usercenter/internal/v1/form"
	"github.com/pengcainiao/zero/rest/httprouter"
)

type User struct {
}

func (u User) RegisterAndLogin(ctx *httprouter.Context, params form.RegisterAndLoginRequest) (form.RegisterAndLoginResponse, error) {
	var data form.RegisterAndLoginResponse
	return data, nil
}
