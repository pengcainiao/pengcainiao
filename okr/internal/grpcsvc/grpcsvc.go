package grpcsvc

import (
	"context"
	"gitlab.com/a16624741591/zero/rpcx/grpcbase"
	"gitlab.com/a16624741591/zero/rpcx/grpcclient/okr"
)

func NewService() grpcbase.ServerBinding {
	initGetOkr()
	return okr.NewBinding(okr.NewService())
}

func initGetOkr() {
	okr.GetOkrHandler = func(ctx context.Context, req okr.GetOkrRequest) grpcbase.Response {
		var data = okr.GetOkrResponse{
			Name: "iayaiyai",
		}
		return grpcbase.Response{
			Data: data,
		}
	}
}
