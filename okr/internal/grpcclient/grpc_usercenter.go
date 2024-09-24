package grpcclient

import (
	"errors"
	"github.com/pengcainiao2/zero/rest/httprouter"
	"github.com/pengcainiao2/zero/rpcx/grpcbase"
	grpcuc "github.com/pengcainiao2/zero/rpcx/grpcclient/usercenter"
	"log"
)

type UserCenterClient struct {
	client grpcuc.Repository
}

func NewUserCenter() *UserCenterClient {
	c, err := grpcbase.DialClient(grpcbase.ServerAddr(grpcbase.UserCenterSVC))
	if err != nil {
		log.Fatal(err)
	}
	client := c.(grpcuc.Repository)
	return &UserCenterClient{client}
}

func (user UserCenterClient) GetUser(ctx *httprouter.Context, params grpcuc.GetUserRequest) (string, error) {

	resp := user.client.GetUser(ctx, params)
	if resp.Message != "" {
		return "", errors.New(resp.Message)
	}
	r := resp.Data.(grpcuc.GetUserResponse)
	return r.Name, nil
}
