package tracer

import (
	"context"
	r "github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/exporters/trace/zipkin"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"os"
	"runtime"
	"sapi/pkg/client/redis"
	"sapi/pkg/logger"
	"sapi/pkg/model"
	"sapi/pkg/tracer/hook"
	"sync"
)

var (
	m sync.RWMutex
	isInit bool
)

type Option func(*Options)

func Init(option *Options) (err error) {
	m.Lock()
	defer m.Unlock()

	if isInit {
		logger.Warn("已经初始化过Tracer")
		return
	}

	if option.Driver == "jaeger" {
		initJaeger(option)
	} else if option.Driver == "zipkin" {
		initZipkin(option)
	} else {
		initStdout()
	}

	isInit = true
	return
}

func AddHookSpanCtx(ctx context.Context) {
	if redis.Client != nil {
		redis.Client.AddHook(hook.RedisHook{})
	}

	if redis.ClusterClient != nil {
		redis.ClusterClient.ForEachShard(ctx, func(ctx context.Context, shard *r.Client) error {
			shard.AddHook(hook.RedisHook{})
			return nil
		})
	}

	if model.DB != nil {
		dbHook := hook.DbHook{
			Ctx:ctx,
		}
		model.DB.Callback().Query().Replace("gorm:after_query", dbHook.AfterFind)
		model.DB.Callback().Create().Replace("gorm:after_create", dbHook.AfterCreate)
		model.DB.Callback().Update().Replace("gorm:after_update", dbHook.AfterUpdate)
		model.DB.Callback().Delete().Replace("gorm:after_delete", dbHook.AfterDelete)
	}

}

func initStdout() {
	pusher, err := stdout.InstallNewPipeline([]stdout.Option{
		stdout.WithQuantiles([]float64{0.5, 0.9, 0.99}),
		stdout.WithPrettyPrint(),
	}, nil)
	if err != nil {
		log.Fatalf("failed to initialize stdout export pipeline: %v", err)
	}
	defer pusher.Stop()
}

func initJaeger(option *Options) {
	endpointOption := jaeger.WithCollectorEndpoint(option.Jaeger.EndpointUrl)
	if option.Jaeger.Mode == "local" {
		endpointOption = jaeger.WithAgentEndpoint(option.Jaeger.AgentEndpoint)
	}

	hostname , err:= os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	flush, err := jaeger.InstallNewPipeline(
		endpointOption,
		jaeger.WithProcess(jaeger.Process{
			ServiceName: option.Name,
			Tags: []label.KeyValue{
				label.String("hostname", hostname),
				label.String("version", runtime.Version()),
			},
		}),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer flush()
}

func initZipkin(option *Options) {
	var logger = log.New(os.Stderr, option.Name, log.Ldate|log.Ltime|log.Llongfile)
	err := zipkin.InstallNewPipeline(
		option.Zipkin.EndpointUrl,
		option.Name,
		zipkin.WithLogger(logger),
		zipkin.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		log.Fatal(err)
	}
}