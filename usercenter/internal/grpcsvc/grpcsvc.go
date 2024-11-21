package grpcsvc

import (
	"context"
	"gitlab.com/a16624741591/zero/rpcx/grpcbase"
	"gitlab.com/a16624741591/zero/rpcx/grpcclient/usercenter"
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
