package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pengcainiao/pengcainiao/okr/internal/v1/services"
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
// @Success 200 {object} httprouter.Response
// @Failure 400 {object} httprouter.Response
// @Router /v2/task [post]
func (o ObjectiveController) First(c *gin.Context) {
	fmt.Println("11111111")
}
