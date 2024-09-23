package grpcsvc

import (
	"context"
	"github.com/pengcainiao2/zero/rpcx/grpcbase"
	"github.com/pengcainiao2/zero/rpcx/grpcclient/usercenter"
)

func NewService() grpcbase.ServerBinding {
	initGetUsercenter()
	return usercenter.NewBinding(usercenter.NewService())
}

func initGetUsercenter() {
	usercenter.GetUserHandler = func(ctx context.Context, req usercenter.GetUserRequest) grpcbase.Response {
		var data = usercenter.GetUserResponse{
			Name: "ososo",
		}
		return grpcbase.Response{
			Data: data,
		}
	}
}
