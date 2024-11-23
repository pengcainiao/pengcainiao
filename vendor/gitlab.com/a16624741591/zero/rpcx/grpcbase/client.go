package grpcbase

import (
	"context"
	"fmt"
	"os"
	"sync"

	"gitlab.com/a16624741591/zero/rpcx/grpcbase/pool"
)

var (
	clients     = make(map[string]ClientBinding)
	rpcConnPool = sync.Map{} //make(map[string]pool.Pool)
)

// RegisterClients 注册grpc客户端
func RegisterClients(serviceName string, binding ClientBinding) {
	clients[serviceName] = binding
}

func DialClient(serviceName string) (interface{}, error) {
	var (
		srvName  string
		err      error
		connPool pool.Pool
	)
	clientPem := pool.ClientPem()
	if len(clientPem) != 0 {
		var domainName = "flyele.vip"
		if os.Getenv("RELEASE_MODE") == "production" {
			domainName = "flyele.net"
		}
		srvName = fmt.Sprintf("%s.%s:443", serviceName, domainName)
	} else {
		srvName = fmt.Sprintf("%s:%d", serviceName, GRPCServerPort)
	}
	if p, ok := rpcConnPool.Load(serviceName); !ok || p == nil {
		connPool, err = pool.New(srvName, pool.DefaultOptions)
		if err != nil {
			println(err.Error())
		}
		rpcConnPool.Store(serviceName, connPool)
	} else {
		connPool = p.(pool.Pool)
	}

	conn, _ := connPool.Get()
	return clients[serviceName].GRPCClient(conn.Value()), nil
}

// InjectContext  注入context
func InjectContext(ctx context.Context) context.Context {
	newCtx := context.Background()
	if token := ctx.Value(RequestAuthorization); token != nil {
		newCtx = context.WithValue(newCtx, RequestAuthorization, token.(string)) // nolint
	}
	if userID := ctx.Value(RequestUserID); userID != nil {
		newCtx = context.WithValue(newCtx, RequestUserID, userID.(string)) // nolint
	}
	if platform := ctx.Value(RequestPlatform); platform != nil {
		newCtx = context.WithValue(newCtx, RequestPlatform, platform.(string)) // nolint
	}
	return newCtx
}
