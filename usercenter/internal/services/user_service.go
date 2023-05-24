package services

import (
	"errors"
	"github.com/pengcainiao/pengcainiao/usercenter/internal/v1/form"
	"github.com/pengcainiao/pengcainiao/usercenter/internal/v1/models"
	"github.com/pengcainiao/zero/core/logx"
	"github.com/pengcainiao/zero/rest/httprouter"
)

type UserServices struct {
	UserID string
}

func (u UserServices) RegisterUser(ctx *httprouter.Context, params form.RegisterUserRequest) httprouter.Response {
	// 检测用户名是否已经被注册
	if (&models.User{}).UserIsExist(ctx, params.UserName) {
		return httprouter.GetError(httprouter.ErrInvalidParameterCode, errors.New("账号已经被注册")).SetRequestParameter(params)
	}

	resp, err := models.User{}.RegisterUser(ctx, params)
	if err != nil {
		logx.NewTraceLogger(ctx).Err(err).Msg("RegisterUser fail")
		return httprouter.GetError(httprouter.ErrInternalErrorCode, err).SetRequestParameter(params)
	}
	return httprouter.Success(map[string]interface{}{
		"data": resp,
	})
}

func (u UserServices) UserLogin(ctx *httprouter.Context, params form.LoginRequest) httprouter.Response {
	// 检测用户名是否存在
	var user = &models.User{
		UserName: params.UserName,
		Password: params.Password,
	}

	if !user.UserIsExist(ctx, params.UserName) {
		return httprouter.GetError(httprouter.ErrInvalidParameterCode, errors.New("账号不存在")).SetRequestParameter(params)
	}

	resp := user.UserLogin(ctx, params)
	return httprouter.Success(map[string]interface{}{
		"data": resp,
	})
}
