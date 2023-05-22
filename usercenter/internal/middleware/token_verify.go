package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pengcainiao/pengcainiao/usercenter/internal/v1/models"
	"github.com/pengcainiao/pengcainiao/usercenter/internal/v1/utils"
	"github.com/pengcainiao/zero/core/logx"
	"github.com/pengcainiao/zero/rest/httprouter"
	"net/url"
	"strings"
)

var (
	TokenVerifyPath = []string{"/v1/user/verify"}
)


func TokenVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			req          = c.Request
			xoriginalURL = getOriginalURL(c)
			authCode     = getUserToken(c, xoriginalURL)

			method = req.Method
			path   = req.URL.Path
		)
		//writeCookie(c)
		if xoriginalURL != nil && xoriginalURL.Path != "" {
			method = c.Request.Header.Get("X-Original-Method")
			path = xoriginalURL.Path
		}

		if authCode != "" {
			token := models.VerifyToken(authCode)
			if token == nil {
				if tireRoute.Handle(method, path) || strings.Contains(path, "auth/verify_code") {
					c.Request.Header.Add("auth_ignore", "1")
					c.Next()
					return
				}
				c.AbortWithStatusJSON(401, gin.H{
					"code":    httprouter.ErrUnAuthorizedCode,
					"message": "token已失效，请重新登录",
				})
				return
			} else {

				// 临时注销，黑名单配置在etcd
				//if !checkHasAuth(c, token.UserID) {
				//	return
				//}
				if !checkUserCanRequestRelease(c, token.UserID) {
					return
				}
				injectResponseHeader(c, token)
				c.Next()
				return
			}
		} else if tireRoute.Handle(method, path) || strings.Contains(path, "auth/verify_code") {
			// 如果未设置token则检查是否不要求
			c.Request.Header.Add("auth_ignore", "1")
			c.Next()
			return
		}
		logx.NewTraceLogger(c.Request.Context()).Info().Interface("req_header", c.Request.Header).
			Str("path", path).
			Str("req_method", req.Method).Msg("需要鉴权")
		c.AbortWithStatusJSON(401, gin.H{
			"code":    httprouter.ErrUnAuthorizedCode,
			"message": "不合法的token参数",
		})

	}
}

func checkUserCanRequestRelease(c *gin.Context, uID string) bool {
	releaseCheck := models.NewUserReleaseCheck()
	if canReq := releaseCheck.CheckCanReq(uID); !canReq {
		logx.NewTraceLogger(c.Request.Context()).Info().Str("userId", uID).Msg("无权限访问此环境")
		c.AbortWithStatusJSON(401, gin.H{
			"code":    httprouter.ErrReqReleaseAuthCode,
			"message": "无权限访问此环境",
		})
		return false
	}
	return true
}

func injectResponseHeader(c *gin.Context, token *models.TokenObject) {
	if token == nil {
		return
	}
	path := c.Request.URL.Path
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), utils.TokenContext, token))

	c.Request.Header.Add("X-Auth-Platform", string(token.Platform))
	c.Request.Header.Add("X-AUTH-USER", token.UserID)
	c.Request.Header.Add("X-AUTH-Version", token.ClientVersion)
	c.Request.Header.Add("X-AUTH-Device", token.DeviceID)

	for _, tokenPath := range TokenVerifyPath {
		if tokenPath == path {
			// 只有路径为 /user/verify的才返回相关参数
			c.Writer.Header().Add("X-Auth-Platform", string(token.Platform))
			c.Writer.Header().Add("X-AUTH-USER", token.UserID)
			c.Writer.Header().Add("X-AUTH-Version", token.ClientVersion)
			c.Writer.Header().Add("X-AUTH-Device", token.DeviceID)
			c.Writer.Header().Add("X-Request-Id", token.DeviceID)
			break
		}
	}
}

func checkHasAuth(c *gin.Context, uID string) bool {
	var ub = models.UserBlackList{}
	if isBlack := ub.IsBlackUser(c, uID); isBlack {
		logx.NewTraceLogger(c.Request.Context()).Info().Str("userId", uID).Msg("黑名单，无权限访问")
		c.AbortWithStatusJSON(401, gin.H{
			"code":    httprouter.ErrUnAuthorizedCode,
			"message": "无权限访问",
		})
		return false
	}
	return true
}


func getOriginalURL(c *gin.Context) *url.URL {
	xoriginalURL := c.GetHeader("X-Original-Url")
	if xoriginalURL == "" {
		xoriginalURI := c.GetHeader("X-Original-Uri")
		xoriginalURL = xoriginalURI
	}
	u, _ := url.Parse(xoriginalURL)
	return u
}

func getUserToken(c *gin.Context, xoriginalURL *url.URL) string {
	authCode := c.GetHeader("Authorization")
	if authCode == "" && xoriginalURL != nil {
		authCode = xoriginalURL.Query().Get("token")
	}
	return authCode
}

