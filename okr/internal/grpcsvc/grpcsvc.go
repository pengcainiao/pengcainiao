package grpcsvc

import (
	"context"
	"github.com/pengcainiao2/zero/rpcx/grpcbase"
	"github.com/pengcainiao2/zero/rpcx/grpcclient/okr"
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
