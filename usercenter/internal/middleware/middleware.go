package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/martian/log"
	"github.com/pengcainiao2/apierror"
	"github.com/pengcainiao2/usercenter/internal/auth"
	v1 "github.com/pengcainiao2/usercenter/internal/v1"
	"github.com/pengcainiao2/usercenter/internal/v1/models"
	"net/http"
	"strconv"
)

const (
	_uidKey              = "uid"
	_nameKey             = "name"
	_accountKey          = "account"
	_langKey             = "accept-language"
	_authorizationKey    = "Authorization" // 鉴权还没接入 避免完全暴露在公网
	_authorizationBearer = "bearer"
	messengerNs          = "messenger"
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

func AuthenticatedHandlev2() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		AuthenticatedFromGinContext(ctx)
	}
}

type PermissionEnum struct {
}

func AuthenticatedFromGinContext(c *gin.Context) {
	//reqMethod := c.Request.Method
	reqPath := c.Request.URL.Path[24:]
	token := c.Request.Header.Get("token")
	// todo
	tokenInfo, err := auth.New(&auth.Config{}).ValidateToken(c, &auth.ValidateTokenRequest{
		Type:  auth.TokenType_ACCESS_TOKEN,
		Token: token,
	})
	if err != nil {
		log.Errorf("VerifyAccessToken %s: %v", token, err)
		c.AbortWithStatusJSON(200, gin.H{
			"code": apierror.UnauthorizedErr.Code(),
			"msg":  apierror.UnauthorizedErr.Message(),
		})
		log.Errorf("url : %s abort with 401 for invalid token", c.Request.URL.Path)
		return
	}

	if tokenInfo.TokenInfo.Uid == 0 {
		c.AbortWithStatusJSON(200, gin.H{
			"code": apierror.UnauthorizedErr.Code(),
			"msg":  apierror.UnauthorizedErr.Message(),
		})
		log.Errorf("tokenInfo.TokenInfo.Uid == 0,Path:%v", c.Request.URL.Path)
		return
	}
	log.Debugf("uid : %d", tokenInfo.TokenInfo.Uid)

	level, err := models.GetUidRoles(c, strconv.FormatUint(tokenInfo.TokenInfo.Uid, 10))
	if err != nil || level == 0 {
		c.AbortWithStatusJSON(200, gin.H{
			"code": apierror.UnauthorizedErr.Code(),
			"msg":  apierror.UnauthorizedErr.Message(),
		})
		log.Errorf("url : %s abort with 401 for zero uid in token", c.Request.URL.Path)
		return
	}

	if pass, _ := v1.EnforceList(strconv.FormatInt(level, 10), reqPath); !pass {
		c.AbortWithStatusJSON(200, gin.H{
			"code": apierror.UnauthorizedErr.Code(),
			"msg":  apierror.UnauthorizedErr.Message(),
		})
		log.Errorf("EnforceList Path:%v,level:%v，reqPath:%v", c.Request.URL.Path, level, reqPath)
		return
	}

	c.Set(_uidKey, tokenInfo.TokenInfo.Uid)
	c.Set(_langKey, c.GetHeader(_langKey))
	c.Next()
}
