package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 设置允许跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-OriginType, Origin, Authorization, X-Auth-User")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE, PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-OriginType")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "172800")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		//c.Request.Header.Add("X-AUTH-USER", "1416492717179104") // fixme
		//c.Set("user_id", "1416492717179104")

		// 处理请求
		c.Next()
	}
}
