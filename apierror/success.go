package apierror

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess           = 0
	DefaultSuccessMessage = "ok"
)

// Success 返回规范化成功响应, 通过断言支持更多特性。
//
// 如果参数拥有 Message 方法则会修改默认的 ok 消息，
// 拥有 Data 方法会重新获取 Data 数据，
// 拥有 Extra 方法会赋予响应更多字段。
func Success(c *gin.Context, data interface{}) {
	msg := DefaultSuccessMessage
	if v, ok := data.(interface{ Message() string }); ok {
		msg = v.Message()
	}

	respData := data
	if v, ok := data.(interface{ Data() interface{} }); ok {
		respData = v.Data()
	}

	basicResp := gin.H{"code": CodeSuccess, "msg": msg, "data": respData}

	if v, ok := data.(interface{ Extra() gin.H }); ok {
		for key, value := range v.Extra() {
			basicResp[key] = value
		}
	}
	c.JSON(http.StatusOK, basicResp)
}
