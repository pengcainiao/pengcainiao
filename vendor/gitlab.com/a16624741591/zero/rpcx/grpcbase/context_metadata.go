package grpcbase

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"
)

// GrpcContext grpc需注入的header值
type GrpcContext map[string]interface{}

const (
	RequestAuthorization = "authorization"
	RequestUserID        = "x-auth-user"
	RequestPlatform      = "x-auth-platform"
	RequestClientVersion = "x-auth-version"
	RequestClientIP      = "x-auth-clientip"
	RequestGrpcContext   = "x-rpc-context"
	CorrelationID        = "correlation-id"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func GetGrpcContext(ctx context.Context) GrpcContext {
	var rpcCtx = ctx.Value(RequestGrpcContext)
	if v, ok := rpcCtx.(GrpcContext); ok {
		return v
	}
	return nil
}

func (g GrpcContext) InjectContext(ctx context.Context) context.Context {
	for k, v := range g {
		ctx = context.WithValue(ctx, strings.ToUpper(k), v) //nolint
	}
	return ctx
}

/* client before functions */
func InjectUserID(ctx context.Context, md *metadata.MD) context.Context {
	if hdr := ctx.Value("Authorization"); hdr != nil {
		md.Set(RequestAuthorization, hdr.(string))
	}
	if hdr := ctx.Value("X-Auth-User"); hdr != nil {
		md.Set(RequestUserID, hdr.(string))
	}
	if hdr := ctx.Value("X-Auth-Platform"); hdr != nil {
		md.Set(RequestPlatform, hdr.(string))
	}
	if hdr := ctx.Value("X-Auth-Version"); hdr != nil {
		md.Set(RequestClientVersion, hdr.(string))
	}
	if hdr := ctx.Value("X-Auth-ClientIP"); hdr != nil {
		md.Set(RequestClientIP, hdr.(string))
	}
	if hdr := ctx.Value(CorrelationID); hdr != nil {
		md.Set(CorrelationID, hdr.(string))
	} else {
		md.Set(CorrelationID, randStringBytesMask(7))
	}
	return ctx
}

// ExtractUserID server before functions
// nolint:golint,unused
func ExtractUserID(ctx context.Context, md metadata.MD) context.Context {
	var rpcContext = GrpcContext{}
	for key, m := range md {
		for _, metadatum := range m {
			rpcContext[key] = metadatum
		}
	}
	ctx = context.WithValue(ctx, RequestGrpcContext, rpcContext) //nolint

	if values := md.Get(RequestAuthorization); len(values) > 0 {
		ctx = context.WithValue(ctx, "Authorization", values[0]) //nolint
	}

	if values := md.Get(RequestUserID); len(values) > 0 {
		ctx = context.WithValue(ctx, "X-Auth-User", values[0]) //nolint
	}

	if values := md.Get(RequestPlatform); len(values) > 0 {
		ctx = context.WithValue(ctx, "X-Auth-Platform", values[0]) //nolint
	}

	if values := md.Get(RequestClientVersion); len(values) > 0 {
		ctx = context.WithValue(ctx, "X-Auth-Version", values[0]) //nolint
	}

	if values := md.Get(RequestClientIP); len(values) > 0 {
		ctx = context.WithValue(ctx, "X-Auth-ClientIP", values[0]) //nolint
	}

	return ctx
}

func randStringBytesMask(n int) string {
	var (
		src = rand.NewSource(time.Now().UnixNano())
		sb  = strings.Builder{}
	)
	sb.Grow(n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}
