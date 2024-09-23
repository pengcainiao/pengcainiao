package grpcsvc

import (
	"github.com/pengcainiao2/zero/rpcx/grpcbase"
	usercenter "github.com/pengcainiao2/zero/rpcx/grpcclient/usercenter"
)

func NewService() grpcbase.ServerBinding {
	return usercenter.NewBinding(usercenter.NewService())
}
