package redis

import (
	"strings"
)

type Options struct {
	Mode string `json:"mode"` //single,sentinel,cluster
	Password string `json:"password"`
	DB int `json:"db"`
	Timeout  int `json:"timeout"`
	MaxRetries int `json:"maxRetries"`
	Single  Single `json:"single"`
	Sentinel Sentinel `json:"sentinel"`
	Cluster Cluster `json:"cluster"`
	Pool Pool `json:"pool"`
}

type Single struct {
	Network  string `json:network`
	Addr string `json:addr`
}

type Sentinel struct {
	Master  string `json:"master"`
	SentinelAddrs  string `json:"metric"`
	nodes   []string
}

type Cluster struct {
	ClusterAddrs  string `json:"metric"`
	nodes   []string
}

type Pool struct {
	PoolSize  int `json:"poolSize"`
	MinIdleConns  int `json:"minIdleConns"`
}

func initOptions() Options{
	option := Options{
		Mode:"single",
		Password:"",
		DB:0,
		Timeout:3,
		MaxRetries:3,
		Single: Single{
			Addr:"localhost:6379",
		},
	}

	return option
}

func NewOptions(opts ...Option) Options {
	options := initOptions()
	for _, o := range opts {
		o(&options)
	}

	return options
}

func Mode(mode string) Option {
	return func(o *Options) {
		o.Mode = mode
	}
}

func Password(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

func Db(db int) Option {
	return func(o *Options) {
		o.DB = db
	}
}

func MaxRetries(maxRetries int) Option {
	return func(o *Options) {
		o.MaxRetries = maxRetries
	}
}

func Network(network string) Option {
	return func(o *Options) {
		o.Single.Network = network
	}
}

func Addr(addr string) Option {
	return func(o *Options) {
		o.Single.Addr = addr
	}
}

func Master(master string) Option {
	return func(o *Options) {
		o.Sentinel.Master = master
	}
}

func SentinelAddrs(sentinelAddrs string) Option {
	return func(o *Options) {
		o.Sentinel.SentinelAddrs = sentinelAddrs
		if len(sentinelAddrs) != 0 {
			for _, v := range strings.Split(sentinelAddrs, ",") {
				v = strings.TrimSpace(v)
				o.Sentinel.nodes = append(o.Sentinel.nodes, v)
			}
		}
	}
}

func ClusterAddrs(clusterAddrs string) Option {
	return func(o *Options) {
		o.Cluster.ClusterAddrs = clusterAddrs
		if len(clusterAddrs) != 0 {
			for _, v := range strings.Split(clusterAddrs, ",") {
				v = strings.TrimSpace(v)
				o.Cluster.nodes = append(o.Cluster.nodes, v)
			}
		}
	}
}

func PoolSize(poolSize int) Option {
	return func(o *Options) {
		o.Pool.PoolSize = poolSize
	}
}

func MinIdleConns(minIdleConns int) Option {
	return func(o *Options) {
		o.Pool.MinIdleConns = minIdleConns
	}
}