package conf

import (
	"gitlab.com/a16624741591/zero/core/env"
	"time"
)

type (
	// Option defines the method to customize the config Options.
	Option func(opt *Options)

	Options struct {
		env                    bool
		EtcdAddress            string  //连接地址
		NsqAddress             string  //连接地址
		ZincAddress            string  //zinc连接地址
		JaegerCollectorAddress string  //jaeger连接地址
		JaegerSamplerRate      float64 //jaeger 数据采样率
		MaxInFlight            int     //最大NSQ消费并发数
		MaxNSQGoroutineCount   int     //NSQ用来处理消息的协程数
		RedisPoolSize          int     //Redis连接池
		RedisMinIdleConns      int     //最小连接空闲数
		MySQLMaxOpenConns      int
		MySQLMaxIdleConns      int
		MySQLMaxLifetime       time.Duration
	}
)

// UseEnv customizes the config to use environment variables.
func UseEnv() Option {
	return func(opt *Options) {
		opt.env = true
	}
}

func (c *Options) setDefault() {
	if c.MaxInFlight == 0 {
		c.MaxInFlight = 1000
	}
	if c.MaxNSQGoroutineCount == 0 {
		c.MaxNSQGoroutineCount = 100
	}
	if c.RedisPoolSize == 0 {
		c.RedisPoolSize = 10
	}
	if c.RedisMinIdleConns == 0 {
		c.RedisMinIdleConns = 5
	}
	if c.NsqAddress == "" {
		c.NsqAddress = env.NSQAddress
	}
	if c.EtcdAddress == "" {
		c.EtcdAddress = env.ETCDAddress
	}
	if c.JaegerCollectorAddress == "" {
		c.JaegerCollectorAddress = env.JaegerCollectorAddress
	}
	if c.JaegerSamplerRate == 0 {
		c.JaegerSamplerRate = 1.0
	}
	if c.MySQLMaxOpenConns == 0 {
		c.MySQLMaxOpenConns = 5
	}
	if c.MySQLMaxIdleConns == 0 {
		c.MySQLMaxIdleConns = 5
	}
	if c.MySQLMaxLifetime == 0 {
		c.MySQLMaxLifetime = time.Second * 10
	}
	if c.ZincAddress == "" {
		c.ZincAddress = env.ZincAddress
	}
}

// WithMaxNSQGoroutineCount 设置可以用来处理MQ消息的协程数
func WithMaxNSQGoroutineCount(goroutines int) Option {
	return func(parameter *Options) {
		parameter.MaxNSQGoroutineCount = goroutines
	}
}

func WithMySQLMaxOpenConns(maxOpen int) Option {
	return func(parameter *Options) {
		parameter.MySQLMaxOpenConns = maxOpen
	}
}

func WithMySQLMaxIdleConns(maxIdle int) Option {
	return func(parameter *Options) {
		parameter.MySQLMaxIdleConns = maxIdle
	}
}

func WithMySQLMaxIdleDuration(maxIdleDuration time.Duration) Option {
	return func(parameter *Options) {
		parameter.MySQLMaxLifetime = maxIdleDuration
	}
}

func WithNsqEndpoint(endpoint string) Option {
	return func(parameter *Options) {
		parameter.NsqAddress = endpoint
	}
}

func WithEtcdEndpoint(endpoint string) Option {
	return func(parameter *Options) {
		parameter.EtcdAddress = endpoint
	}
}

func WithJaegerCollectorEndpoint(endpoint string) Option {
	return func(parameter *Options) {
		parameter.JaegerCollectorAddress = endpoint
	}
}
func WithJaegerCollectorSamplerRate(samplerRate float64) Option {
	return func(parameter *Options) {
		parameter.JaegerSamplerRate = samplerRate
	}
}

// WithMaxConcurrency 设置NSQ最大并发数（MaxInFlight），默认值 1000
func WithMaxConcurrency(maxConcurrency int) Option {
	return func(parameter *Options) {
		parameter.MaxInFlight = maxConcurrency
	}
}

func WithRedisPoolSize(poolSize int) Option {
	return func(parameter *Options) {
		parameter.RedisPoolSize = poolSize
	}
}

func WithRedisMinIdleConns(minIdls int) Option {
	return func(parameter *Options) {
		parameter.RedisMinIdleConns = minIdls
	}
}

func ApplyConfig(opts ...Option) *Options {
	var lcfg = &Options{}
	for _, opt := range opts {
		opt(lcfg)
	}
	lcfg.setDefault()
	return lcfg
}
