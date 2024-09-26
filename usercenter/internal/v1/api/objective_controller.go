package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/martian/log"
	"github.com/pengcainiao2/apierror"
	"github.com/pengcainiao2/usercenter/internal/auth"
	"github.com/pengcainiao2/usercenter/internal/v1/form"
	"github.com/pengcainiao2/usercenter/internal/v1/models"
	"github.com/pengcainiao2/usercenter/internal/v1/services"
	"github.com/pengcainiao2/zero/core/logx"
	"github.com/pengcainiao2/zero/rest/httprouter"
	"github.com/pengcainiao2/zero/rpcx/grpcclient/okr"
	"github.com/pengcainiao2/zero/tools/syncer"
	"time"
)

type ObjectiveController struct {
	service services.ObjectiveServices
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
func (o ObjectiveController) First(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		httprouter.ResponseJSONContent(c, httprouter.Success(map[string]interface{}{
			"data": "空值",
		}))
	}
	httprouter.ResponseJSONContent(c, httprouter.Success(map[string]interface{}{
		"data": keyword,
	}))
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
func (o ObjectiveController) GongZhu(c *gin.Context) {
	httprouter.ResponseJSONContent(c, httprouter.Success(map[string]interface{}{
		"data": "https://oversea-test-666.oss-ap-southeast-1.aliyuncs.com/feishu/FkibeHZXlBTfRQV9I29YSbppvEHd.jpeg",
	}))
}

func (o ObjectiveController) TestRedis(c *gin.Context) {
	res, err := syncer.Redis().Get(c, "test").Result()
	if err != nil {
		logx.NewTraceLogger(c).Err(err).Msg("TestRedis fail")
		return
	}

	httprouter.ResponseJSONContent(c, httprouter.Success(map[string]interface{}{
		"data": res,
	}))
}

type Okr struct {
	Id   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (o ObjectiveController) Mysql(c *gin.Context) {
	var (
		Ob   Okr
		asql = `select id,name from okr.okr where id = 1`
	)
	err := syncer.MySQL().Get(c, &Ob, asql)
	if err != nil {
		logx.NewTraceLogger(c).Err(err).Msg("tests Mysql fail")
		return
	}

	httprouter.ResponseJSONContent(c, httprouter.Success(map[string]interface{}{
		"data": Ob,
	}))
}

func (o ObjectiveController) Rpc(c *gin.Context) {
	newClient := okr.NewClient()
	params := okr.GetOkrRequest{}
	ctx := &httprouter.Context{}
	resp := newClient.HandleGetOkr(ctx, params)

	httprouter.ResponseJSONContent(c, httprouter.Success(map[string]interface{}{
		"data": resp.Name,
	}))
}

func (s *ObjectiveController) Login(ctx *gin.Context) {
	var (
		in  form.LoginReq
		out form.LoginResp
	)
	if err := ctx.BindJSON(&in); err != nil {
		apierror.Fail(ctx, apierror.InvalidParamErr)
		return
	}
	logx.NewTraceLogger(ctx).Info().Msg("111")
	key := fmt.Sprintf("coin_agent_login_%v", in.UserName)
	if syncer.Redis().Get(ctx, key).Val() == "0" {
		log.Infof("Get(key).Val() == 0 in:%v", in)
		apierror.Fail(ctx, apierror.LoginTimesOneHourLimitFErr)
		return
	}
	uid, err := models.AuthenticateCoinAgentUser(ctx, in)

	logx.NewTraceLogger(ctx).Info().Msg(fmt.Sprintf("2222,%v", uid))

	if err != nil || uid == 0 {
		logx.NewTraceLogger(ctx).Info().Msg(fmt.Sprintf("login fail in:%v", in))
		if !syncer.Redis().SetNX(ctx, key, 5, 1*time.Hour).Val() {
			if syncer.Redis().TTL(ctx, key).Val().Seconds() > 0 {
				if syncer.Redis().Decr(ctx, key).Val() <= 0 {
					apierror.Fail(ctx, apierror.LoginTimesOneHourLimitFErr)
					return
				}
			}
		}
		apierror.Fail(ctx, apierror.InvalidUserErr)
		return
	}

	resp, err := models.GetAgentUsersInfo(ctx, []uint64{uid})
	if err != nil {
		log.Errorf("login GetAgentUsersInfo fail err:%v", err)
		apierror.Fail(ctx, apierror.InternalErr)
		return
	}
	logx.NewTraceLogger(ctx).Info().Msg(fmt.Sprintf("3333,%v", uid))

	if detail, ok := resp[uid]; ok {
		if detail.Banned == 1 {
			logx.NewTraceLogger(ctx).Info().Msg(fmt.Sprintf("login 賬號被凍結 uid:%v", uid))
			apierror.Fail(ctx, apierror.InvalidUser)
			return
		}
	}
	oauthToken := &auth.TokenInfo{Uid: uid, DeviceId: in.DeviceID}
	if auth.New(&auth.Config{}) == nil {
		logx.NewTraceLogger(ctx).Info().Msg(fmt.Sprintf("444"))

	}

	tokens, err := auth.New(&auth.Config{}).GrantTokens(ctx, &auth.GrantTokensRequest{TokenInfo: oauthToken})
	out.RefreshToken = tokens.RefreshToken.Token
	out.RefreshTokenExpiresIn = tokens.RefreshToken.ExpiresIn
	out.AssessToken = tokens.AccessToken.Token
	out.AssessTokenExpiresIn = tokens.AccessToken.ExpiresIn
	out.Uid = uid
	apierror.Success(ctx, out)
}
