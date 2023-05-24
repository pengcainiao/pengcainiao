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

// RegisterUser
// @Summary 注册
// @Description 注册
// @Tags user（用户管理）
// @accept json
// @Produce json
// @Param   body            body        form.RegisterUserRequest        true    "json数据"
// @Success 200 {object} httprouter.Response  "用户对象"
// @Failure 400 {object} httprouter.Response
// @Router /v2/user/register [post]
func (u UserController) RegisterUser(c *gin.Context) {
	var params form.RegisterUserRequest
	if err := c.ShouldBind(&params); err != nil {
		httprouter.ResponseJSONContent(c, httprouter.ErrorSame(
			httprouter.ErrInvalidParameterCode,
			validator.NewValidator().Translate(err),
		))
		return
	}

	u.service.UserID = c.GetString("user_id")
	res := u.service.RegisterUser(httprouter.NewContext(c), params)
	httprouter.ResponseJSONContent(c, res)
}

// UserLogin
// @Summary 登陆
// @Description 登陆
// @Tags user（用户管理）
// @accept json
// @Produce json
// @Param   body            body        form.LoginRequest        true    "json数据"
// @Success 200 {object} httprouter.Response  "用户对象"
// @Failure 400 {object} httprouter.Response
// @Router /v2/user/login [post]
func (u UserController) UserLogin(c *gin.Context) {
	var params form.LoginRequest
	if err := c.ShouldBind(&params); err != nil {
		httprouter.ResponseJSONContent(c, httprouter.ErrorSame(
			httprouter.ErrInvalidParameterCode,
			validator.NewValidator().Translate(err),
		))
		return
	}

	u.service.UserID = c.GetString("user_id")
	res := u.service.UserLogin(httprouter.NewContext(c), params)
	httprouter.ResponseJSONContent(c, res)
}
