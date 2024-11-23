package grpcbase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"gitlab.com/a16624741591/zero/core/env"
	"gitlab.com/a16624741591/zero/core/trace"
	"gitlab.com/a16624741591/zero/rpcx/grpcbase/pool"
	pb "gitlab.com/a16624741591/zero/rpcx/protos"
	serverinterceptors "gitlab.com/a16624741591/zero/rpcx/serviceinterceptors"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GRPCServerPort grpc服务端口
const (
	UserCenterSVC = "usercenter-svc"
	OkrSVC        = "okr-svc"
)

var (
	GRPCServerPort          = 8084
	PrivateDaemonSVC        = "127.0.0.1:50051"
	PrivateDeployServiceSVC = ""
)

var (
	NotImplErrorResponse = ErrorResponse(errors.New("未实现"))
	kaep                 = keepalive.EnforcementPolicy{
		MinTime:             3 * time.Minute, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}
	kasp = keepalive.ServerParameters{
		//MaxConnectionIdle: 10 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
	}
)

type rateLimiter struct {
	*rate.Limiter
}

func (r rateLimiter) Limit() bool {
	if err := r.Limiter.Wait(context.Background()); err != nil {
		return false
	}
	return r.Limiter.Allow()
}

func newGrpcServer() *grpc.Server {
	//limiter := rateLimiter{rate.NewLimiter(rate.Every(time.Millisecond*5), 1)}
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		serverinterceptors.UnaryCrashInterceptor,
		//serverinterceptors.UnaryStatInterceptor(s.metrics),
		serverinterceptors.UnaryPrometheusInterceptor,
	}
	if env.EnableTracing() {
		unaryInterceptors = append(unaryInterceptors, serverinterceptors.UnaryTracingInterceptor)
	}
	srv := grpc.NewServer(
		grpc.InitialWindowSize(pool.InitialWindowSize),
		grpc.InitialConnWindowSize(pool.InitialConnWindowSize),
		grpc.MaxSendMsgSize(pool.MaxSendMsgSize),
		grpc.MaxRecvMsgSize(pool.MaxRecvMsgSize),
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(opts...),
		),
	)
	return srv
}

// RegisterServer 注册grpc服务
func RegisterServer(binding ServerBinding) error {
	srv := newGrpcServer()

	sc, err := net.Listen("tcp", fmt.Sprintf(":%d", GRPCServerPort))
	if err != nil {
		log.Fatalf("无法创建GRPC server，unable to listen: %+v", err)
	}
	go func() {
		log.Println("已建立GPRC服务")
		if err := srv.Serve(sc); err != nil {
			log.Fatalf("无法创建GRPC server：%v", err)
		}
	}()
	return binding.RegisterServer(srv)
}

// ServerAddr 服务器地址
func ServerAddr(serviceName string) string {
	if env.IsDevMode() {
		return "0.0.0.0"
	}
	return serviceName
}

// CreateGRPCServer 创建GRPC服务
func CreateGRPCServer(endpoint endpoint.Endpoint, decReq grpctransport.DecodeRequestFunc, decResp grpctransport.EncodeResponseFunc) *grpctransport.Server {
	return grpctransport.NewServer(
		endpoint,
		decReq,
		decResp,
	)
}

func CreateGRPCClientEndpoint(conn *grpc.ClientConn, serviceName, method string, enc grpctransport.EncodeRequestFunc, dec grpctransport.DecodeResponseFunc) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		serviceName,
		method,
		enc,
		dec,
		&pb.Response{},
		grpctransport.ClientBefore(func(ctx context.Context, md *metadata.MD) context.Context {
			if v, ok := ctx.Value(trace.TraceIdKey).(string); ok {
				md.Set(trace.TraceIdKey, v)
				ctx = metadata.NewIncomingContext(ctx, *md)
			}
			return ctx
		}),
	).Endpoint()
}
