package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/a16624741591/zero/core/logx"
	"gitlab.com/a16624741591/zero/rest/httprouter"
	grpcuc "gitlab.com/a16624741591/zero/rpcx/grpcclient/usercenter"
	"gitlab.com/a16624741591/zero/tools/syncer"
	"pp/okr/internal/v1/services"
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
			"data": "空值1",
		}))
	}
	httprouter.ResponseJSONContent(c, httprouter.Success(map[string]interface{}{
		"data": keyword + "1",
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

func (o ObjectiveController) GrpcTry(c *gin.Context) {
	ctx := &httprouter.Context{}
	params := grpcuc.GetUserRequest{
		Keyword: "AA",
		Context: &grpcuc.UserContext{
			UserID:        "1",
			Platform:      "1",
			ClientVersion: "1",
			Token:         "1",
			ClientIP:      "1",
			RequestID:     "1",
		},
	}
	newClient := grpcuc.NewClient()
	resp := newClient.HandleGetUser(ctx, params)
	httprouter.ResponseJSONContent(c, httprouter.Success(map[string]interface{}{
		"data": resp.Name,
	}))
}
