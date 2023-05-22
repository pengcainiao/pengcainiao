package services

import (
	"github.com/pengcainiao/pengcainiao/usercenter/internal/v1/form"
	"github.com/pengcainiao/pengcainiao/usercenter/internal/v1/models"
	"github.com/pengcainiao/zero/rest/httprouter"
)

type UserServices struct {

}

func (u UserServices)RegisterAndLogin(ctx *httprouter.Context,params form.RegisterAndLoginRequest)httprouter.Response  {
	resp,err:=models.User{}.RegisterAndLogin(ctx,params)
	if err!=nil{
		return httprouter.Response{}
	}
	return httprouter.Success(map[string]interface{}{
		"data":resp,
	})
}
