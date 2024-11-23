package grpcbase

import (
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

// ServerBinding 服务端绑定
type ServerBinding interface {
	//RegisterServer 注册服务
	RegisterServer(srv *grpc.Server) error
	//GRPCHandler 获取grpc处理程序
	GRPCHandler() map[string]grpctransport.Handler
}

// ClientBinding 客户端绑定
type ClientBinding interface {
	//GRPCClient 获取客户端接口类型
	GRPCClient(cc *grpc.ClientConn) interface{}
}
