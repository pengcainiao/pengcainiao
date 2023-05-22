package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pengcainiao/pengcainiao/usercenter/internal/services"
	"github.com/pengcainiao/pengcainiao/usercenter/internal/v1/form"
	"github.com/pengcainiao/zero/rest/httprouter"
	"github.com/pengcainiao/zero/rest/validator"
)

type UserController struct {
	service services.UserServices
}

// First 测试
// @Summary 测试
// @Tags 测试（任务）
// @Security ApiKeyAuth
// @accept json
// @Produce json
// @Param   keyword     query       string                  false   "任务ID, 多个用逗号隔开"
// @Success 200 {object} httprouter.Response
// @Failure 400 {object} httprouter.Response
// @Router /v2/test [get]
func (u UserController) First(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		httprouter.ResponseJSONContent(c, httprouter.Success(map[string]interface{}{
			"data": "user:" + "空值",
		}))
	}
	httprouter.ResponseJSONContent(c, httprouter.Success(map[string]interface{}{
		"data": "user:" + keyword,
	}))
}

// RegisterAndLogin
// @Summary 注册登录
// @Description 注册登录
// @Tags user（用户管理）
// @accept json
// @Produce json
// @Param   body            body        form.RegisterAndLoginRequest        true    "json数据"
// @Success 200 {object} httprouter.Response  "用户对象"
// @Failure 400 {object} httprouter.Response
// @Router /v2/auth/phonelogin [post]
func (u UserController)RegisterAndLogin(c *gin.Context)  {
	var params form.RegisterAndLoginRequest
	if err := c.ShouldBind(&params); err != nil {
		httprouter.ResponseJSONContent(c, httprouter.ErrorSame(
			httprouter.ErrInvalidParameterCode,
			validator.NewValidator().Translate(err),
		))
		return
	}
	res:=u.service.RegisterAndLogin(httprouter.NewContext(c),params)
	httprouter.ResponseJSONContent(c,res)
}
