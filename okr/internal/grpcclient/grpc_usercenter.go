package grpcclient

//
//import (
//	"github.com/pengcainiao2/zero/rest/httprouter"
//	grpcuc "github.com/pengcainiao2/zero/rpcx/grpcclient/usercenter"
//	"log"
//)
//
//type UserCenterClient struct {
//	client grpcuc.Repository
//}
//
////func NewUserCenter() *UserCenterClient {
////	c := grpcuc.NewClient()
////	//c, err := grpcbase.DialClient(grpcbase.ServerAddr(grpcbase.UserCenterSVC))
////	//if err != nil {
////	//	log.Println("NewUserCenter fail")
////	//	log.Fatal(err)
////	//}
////	//client := c.(grpcuc.Repository)
////	//log.Println("--", client)
////	//return &UserCenterClient{client}
////	return c
////}
//
//func (user UserCenterClient) GetUser(ctx *httprouter.Context, params grpcuc.GetUserRequest) (string, error) {
//	log.Println("22222")
//	if user.client == nil {
//		log.Println("123123")
//	} else {
//		log.Println("23233", user.client)
//	}
//	newClient := grpcuc.NewClient()
//	resp := newClient.HandleGetUser(ctx, params)
//
//	r := resp.(grpcuc.GetUserResponse)
//	return r.Name, nil
//}
