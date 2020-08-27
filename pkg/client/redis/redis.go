package redis

import (
	"context"
	r "github.com/go-redis/redis/v8"
	"sapi/pkg/logger"
	"sync"
	"time"
)

var (
	Client *r.Client
	ClusterClient *r.ClusterClient
	options *Options
	m sync.RWMutex
	isInit bool
)

type Option func(*Options)

func Init(ctx context.Context, option *Options) (err error) {
	m.Lock()
	defer m.Unlock()

	if isInit {
		logger.Warn("已经初始化过Redis")
		return
	}

	logger.Debug(option)
	options = option

	if option.Mode == "single" {
		Client = initSingle(option)
	} else if option.Mode == "sentinel" {
		Client = initSentinel(option)
	} else {
		ClusterClient = initCluster(option)
	}

	if Client != nil {
		_, err = Client.Ping(ctx).Result()
		if err == nil {
			isInit = true
			logger.Info("初始化redis完成")
		} else {
			logger.Error(err)
		}
	}

	if ClusterClient != nil {
		err = ClusterClient.ForEachShard(ctx, func(ctx context.Context, shard *r.Client) error {
			return shard.Ping(ctx).Err()
		})
		if err == nil {
			isInit = true
			logger.Info("redis cluster初始化成功")
		} else {
			logger.Error(err)
		}
	}

	return err
}

func initSingle(option *Options) *r.Client{
	return r.NewClient(&r.Options{
		Network: option.Single.Network,
		Addr:     option.Single.Addr,
		DB:       option.DB,
		Password:      option.Password,
		MaxRetries:    option.MaxRetries,
		ReadTimeout:   time.Duration(option.Timeout) * time.Second,
		WriteTimeout:  time.Duration(option.Timeout) * time.Second,
		PoolSize:   option.Pool.PoolSize,
		MinIdleConns: option.Pool.MinIdleConns,
	})
}

func initSentinel(option *Options) *r.Client{
	return r.NewFailoverClient(&r.FailoverOptions{
		MasterName:    option.Sentinel.Master,
		SentinelAddrs: option.Sentinel.nodes,
		DB:            option.DB,
		Password:      option.Password,
		MaxRetries:    option.MaxRetries,
		ReadTimeout:   time.Duration(option.Timeout) * time.Second,
		WriteTimeout:  time.Duration(option.Timeout) * time.Second,
		PoolSize:   option.Pool.PoolSize,
		MinIdleConns: option.Pool.MinIdleConns,
	})
}

func initCluster(option *Options) *r.ClusterClient{
	return r.NewClusterClient(&r.ClusterOptions{
		Addrs: option.Cluster.nodes,
		Password:      option.Password,
		MaxRetries:    option.MaxRetries,
		ReadTimeout:   time.Duration(option.Timeout) * time.Second,
		WriteTimeout:  time.Duration(option.Timeout) * time.Second,
		PoolSize:   option.Pool.PoolSize,
		MinIdleConns: option.Pool.MinIdleConns,
	})
}