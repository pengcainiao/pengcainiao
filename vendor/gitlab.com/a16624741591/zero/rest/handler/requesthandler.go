package handler

import (
	"bytes"
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"gitlab.com/a16624741591/zero/core/logx"
	"gitlab.com/a16624741591/zero/rest/httprouter"
)

func Request() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.URL.String(), "/user/verify") {
			c.Next()
			return
		}

		var (
			traceIdKey = http.CanonicalHeaderKey("x-request-id")
			headerData httprouter.HeaderData
		)

		if err := c.ShouldBindHeader(&headerData); err != nil {
			headerData.UserID = os.Getenv("X_AUTH_USER")
			headerData.Platform = os.Getenv("X_AUTH_PLATFORM")
			headerData.ClientVersion = os.Getenv("X_AUTH_VERSION")
		}

		if headerData.UserID == "" && headerData.Authorization != "" {
			if arr := strings.Split(headerData.Authorization, "."); len(arr) == 3 {
				if bytes, err := base64.RawStdEncoding.DecodeString(arr[1]); err == nil {
					headerData.UserID = gjson.GetBytes(bytes, "UserID").String()
					headerData.Platform = gjson.GetBytes(bytes, "Platform").String()
					headerData.ClientVersion = gjson.GetBytes(bytes, "ClientVersion").String()
				}
			}
		}

		if v := c.GetHeader(traceIdKey); v != "" {
			headerData.RequestID = v
		}

		request, _ := c.GetRawData()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(request))

		c.Set("user_id", headerData.UserID)
		c.Set("header", headerData)

		// 打印数据，把读过的字节流重新放到body
		logEvent := logx.NewTraceLogger(context.Background()).Debug().
			Interface("header", headerData).
			Str("method", c.Request.Method).
			Str("url", c.Request.URL.String())

		if len(request) > 0 {
			c.Set("request", request)
			logEvent.RawJSON("request", request)
		}

		logEvent.Msg("请求参数")

		// 处理请求
		c.Next()
	}
}
