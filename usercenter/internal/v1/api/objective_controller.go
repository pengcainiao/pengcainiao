package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pengcainiao2/usercenter/internal/v1/services"
	"github.com/pengcainiao2/zero/core/logx"
	"github.com/pengcainiao2/zero/rest/httprouter"
	"github.com/pengcainiao2/zero/tools/syncer"
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
