package grpcclient

import (
	"github.com/pengcainiao2/zero/rest/httprouter"
	"github.com/pengcainiao2/zero/rpcx/grpcbase"
	grpcusercenter "github.com/pengcainiao2/zero/rpcx/grpcclient/usercenter"
)

// client grpc客户端
func client(serverName string) (interface{}, error) {
	return grpcbase.DialClient(grpcbase.ServerAddr(serverName))
}

// injectContext 注入context
func injectContext(ctx *httprouter.Context, userContext interface{}) {
	var (
		header    = ctx.Data
		token     = header.Authorization
		userID    = header.UserID
		platform  = header.Platform
		version   = header.ClientVersion
		clientIP  = header.ClientIP
		requestID = header.RequestID
	)

	switch userCtx := userContext.(type) {
	case *grpcusercenter.UserContext:
		userCtx.Token = token
		userCtx.UserID = userID
		userCtx.Platform = platform
		userCtx.ClientVersion = version
		userCtx.ClientIP = clientIP
		userCtx.RequestID = requestID
	}
}
