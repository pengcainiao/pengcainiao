package trace

import (
	oteltrace "go.opentelemetry.io/otel/trace"
)

// HttpTraceName represents the tracing name.
const (
	HttpTraceName    = "http-server"
	RedisTraceName   = "redis-client"
	MysqlTraceName   = "mysql-client"
	ElasticTraceName = "es-client"
	RpcTraceName     = "grpc"
	NsqTraceName     = "nsq"
	ZincTraceName    = "zinc"
)

// A Config is a opentelemetry config.
type Config struct {
	Name     string  `json:",optional"`
	Endpoint string  `json:",optional"`
	Sampler  float64 `json:",default=1.0"`
	Batcher  string  `json:",default=jaeger,options=jaeger|zipkin"`
}

func PatchSpanData(sp oteltrace.Span) {
	if sp == nil {
		return
	}
	sp.SetAttributes()
}
