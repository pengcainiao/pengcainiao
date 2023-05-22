package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pengcainiao/pengcainiao/usercenter/internal/services"
	"github.com/pengcainiao/zero/rest/httprouter"
)

type AccountController struct {
	svc    services.AccountService
	//wechat services.WechatAccountService
	user   services.UserService
}

func NewAccountController() AccountController {
	return AccountController{
		svc:    services.AccountService{},
		//wechat: services.NewWechatAccountService(),
	}
}

//VerifyAccessibleHandler 验证可访问性
// @Summary traefik token验证
// @Description traefik forward_auth验证，返回2XX则表示通过，且会在Header中加入 x-auth-user（用户ID）
// @Accept plain
// @Security ApiKeyAuth
// @Tags V2account（账号）
// @Router /v1/user/verify [get]
func (ac AccountController) VerifyAccessibleHandler(c *gin.Context) {
	httprouter.NewContext(c)
	httprouter.ResponseJSONContent(c, gin.H{})
}