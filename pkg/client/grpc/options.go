package grpc

import (
	"google.golang.org/grpc"
	"time"
)

type Options struct {
	PoolSize int
	PoolTTL  time.Duration
	PoolMaxIdle int
	PoolMaxStreams	int

	grpcDialOptions []grpc.DialOption
	address string
}

type Option func(*Options)

func NewOptions(opts ...Option) *Options {
	options := &Options{}
	for _, o := range opts {
		o(options)
	}

	return options
}

func Address(address string) Option {
	return func(o *Options) {
		o.address = address
	}
}

func GrpcDialOptions(opts ...grpc.DialOption) Option {
	return func(o *Options) {
		o.grpcDialOptions = opts
	}
}

func (opt *Options) Build() *Client {
	if opt.PoolSize == 0 {
		opt.PoolSize = 10
	}
	if opt.PoolTTL == 0 {
		opt.PoolTTL = 3
	}
	if opt.PoolMaxIdle == 0 {
		opt.PoolMaxIdle = 50
	}
	if opt.PoolMaxStreams == 0 {
		opt.PoolMaxStreams = 20
	}
	return newClient(opt)
}

