package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pengcainiao/pengcainiao/okr/internal/v1/services"
	"github.com/pengcainiao/zero/rest/httprouter"
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
