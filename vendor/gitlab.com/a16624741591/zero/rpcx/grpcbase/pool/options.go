// Copyright 2019 shimingyah. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// ee the License for the specific language governing permissions and
// limitations under the License.

package pool

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"time"

	"gitlab.com/a16624741591/zero/core/env"

	"gitlab.com/a16624741591/zero/rpcx/clientinterceptors"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	// DialTimeout the timeout of create connection
	DialTimeout = 5 * time.Second

	// BackoffMaxDelay provided maximum delay when backing off after failed connection attempts.
	BackoffMaxDelay = 3 * time.Second

	// KeepAliveTime is the duration of time after which if the client doesn't see
	// any activity it pings the server to see if the transport is still alive.
	KeepAliveTime = time.Duration(10) * time.Second

	// KeepAliveTimeout is the duration of time for which the client waits after having
	// pinged for keepalive check and if no activity is seen even after that the connection
	// is closed.
	KeepAliveTimeout = time.Duration(3) * time.Second

	// InitialWindowSize we set it 1GB is to provide system's throughput.
	InitialWindowSize = 1 << 30

	// InitialConnWindowSize we set it 1GB is to provide system's throughput.
	InitialConnWindowSize = 1 << 30

	// MaxSendMsgSize set max gRPC request message size sent to server.
	// If any request message size is larger than current value, an error will be reported from gRPC.
	MaxSendMsgSize = 4 << 30

	// MaxRecvMsgSize set max gRPC receive message size received from server.
	// If any message size is larger than current value, an error will be reported from gRPC.
	MaxRecvMsgSize = 4 << 30
)

// Options are params for creating grpc connect pool.
type Options struct {
	// Dial is an application supplied function for creating and configuring a connection.
	Dial func(address string) (*grpc.ClientConn, error)

	// Maximum number of idle connections in the pool.
	MaxIdle int

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int

	// MaxConcurrentStreams limit on the number of concurrent streams to each single connection
	MaxConcurrentStreams int

	// If Reuse is true and the pool is at the MaxActive limit, then Get() reuse
	// the connection to return, If Reuse is false and the pool is at the MaxActive limit,
	// create a one-time connection to return.
	Reuse bool
}

// DefaultOptions sets a list of recommended options for good performance.
// Feel free to modify these to suit your needs.
var DefaultOptions = Options{
	Dial:                 Dial,
	MaxIdle:              8,
	MaxActive:            12,
	MaxConcurrentStreams: 12,
	Reuse:                true,
}

// Dial return a grpc connection with defined configurations.
func Dial(address string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DialTimeout)
	defer cancel()

	var unaryInterceptors = make([]grpc.UnaryClientInterceptor, 0)
	unaryInterceptors = append(unaryInterceptors, clientinterceptors.PrometheusInterceptor)
	if env.EnableTracing() {
		unaryInterceptors = append(unaryInterceptors, clientinterceptors.UnaryTracingInterceptor)
	}
	return grpc.DialContext(ctx, address,
		getTls(address),
		grpc.WithChainUnaryInterceptor(
			unaryInterceptors...,
		),
		grpc.WithStreamInterceptor(
			grpc_retry.StreamClientInterceptor(
				grpc_retry.WithPerRetryTimeout(1*time.Second),
				grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100*time.Millisecond)),
				grpc_retry.WithCodes(codes.NotFound, codes.Aborted))),
		//grpc.WithBalancerName(roundrobin.Name),
		grpc.WithBackoffMaxDelay(BackoffMaxDelay),
		grpc.WithInitialWindowSize(InitialWindowSize),
		grpc.WithInitialConnWindowSize(InitialConnWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(MaxSendMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                KeepAliveTime,
			Timeout:             KeepAliveTimeout,
			PermitWithoutStream: true,
		}))
}

func getTls(address string) grpc.DialOption {
	var (
		clientPem       = ClientPem()
		transCredential credentials.TransportCredentials
	)
	if len(clientPem) > 0 {
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(clientPem) {
			return nil
		}
		// InsecureSkipVerify 跳过证书验证
		transCredential = credentials.NewTLS(&tls.Config{ServerName: address, RootCAs: cp, InsecureSkipVerify: true})
		return grpc.WithTransportCredentials(transCredential)
	}
	return grpc.WithInsecure()
}

// DialTest return a simple grpc connection with defined configurations.
func DialTest(address string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DialTimeout)
	defer cancel()
	return grpc.DialContext(ctx, address, grpc.WithInsecure())
}

func ClientPem() []byte {
	if env.IsDevMode() {
		return nil
	}
	var pem = os.Getenv("GRPCTLS_PEM")
	if pem == "" {
		return []byte("") //defaultClientPem
	}
	return []byte(pem)
}
