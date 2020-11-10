package boot

import (
	"github.com/26597925/EastCloud/cmd/hello/boot/engine"
	"github.com/26597925/EastCloud/cmd/hello/router"
	"github.com/26597925/EastCloud/internal/agent"
	"github.com/26597925/EastCloud/pkg/bootstrap"
	"github.com/26597925/EastCloud/pkg/bootstrap/flag"
	"github.com/26597925/EastCloud/pkg/client/etcdv3"
	"github.com/26597925/EastCloud/pkg/client/redis"
	capi "github.com/26597925/EastCloud/pkg/config/api"
	"github.com/26597925/EastCloud/pkg/logger"
	"github.com/26597925/EastCloud/pkg/model"
	"github.com/26597925/EastCloud/pkg/registry"
	"github.com/26597925/EastCloud/pkg/registry/etcd"
	"github.com/26597925/EastCloud/pkg/server/api"
	"github.com/26597925/EastCloud/pkg/tracer"
	"github.com/google/uuid"
	"math"
	"time"
)

func Init (flags ...flag.Flag) bootstrap.EngineContext {
	fg := []flag.Flag{&flag.StringFlag{
		Name:    "config",
		Usage:   "--config",
		EnvVar:  "CONFIG_PATH",
		Default: "config/config.yml",
	}}
	engine.SetFlag(append(fg, flags...))

	return engine.GetServiceContext()
}

func InitFlag () error {
	fs := flag.NewFlagSet()
	fs.Register(engine.GetFlag() ...)
	err := fs.Parse()
	if err == nil {
		engine.SetFlagSet(fs)
	}
	return err
}

func InitConfig() error {
	cf, err := engine.ParseConfig(engine.GetFlagSet())
	if err != nil {
		return err
	}

	go func() {
		for {
			var data interface{}
			data, err = cf.Watcher.Next()
			if err == nil {
				b, err := capi.Encoders["json"].Encode(data)
				if err != nil {
					return
				}

				capi.Encoders["json"].Decode(b, cf)
			}
		}
	}()
	engine.SetConfig(cf)
	return err
}

func InitLog() error {
	opts := logger.NewOptions(logger.Merge(engine.GetConfig().Logger))
	log, err := logger.NewZap(opts)

	if err != nil {
		return err
	}

	logger.SetLog(log)
	return nil
}

func InitRedis() error {
	return redis.Init(engine.GetContext(), engine.GetConfig().Redis)
}

func InitModel() error {
	return model.Init(engine.GetConfig().Orm)
}

func InitTracer() error {
	tracer.AddHookSpanCtx(engine.GetContext())
	err := tracer.Init(engine.GetConfig().Tracer)
	return err
}

func InitServer() error {
	 handlers := map[string]api.Handler{
	 	"grpc": router.GrpcRouter,
	 	"http": router.HttpRouter,
	 }

	for _, svrOpt := range engine.GetConfig().Server {
		if svrOpt.Id == "" {
			svrOpt.Id = uuid.Must(uuid.NewRandom()).String()
		}

		svr := server.NewServer(svrOpt)
		err := svr.Init()
		if err != nil {
			return err
		}

		svr.Handler(handlers[svr.GetOption().GetName()])
		engine.AddServe(svr)
	}
	return nil
}

func InitRegistry() error {
	var ttl int64
	ttl = 6
	cli := etcdv3.NewOptions().Build()
	opt := &registry.Options{Timeout:3, TTL: ttl}
	rsy, err := etcd.NewRegistry(opt, cli)
	if err != nil {
		return err
	}

	engine.SetRegistry(multi.New(rsy))
	for _, svrOpt := range engine.GetConfig().Server {
		t := time.Duration(math.Ceil(float64(ttl/3))) * time.Second
		err = register(map[string]interface{}{"opt": svrOpt, "ttl": t})
		if err != nil {
			return err
		}
	}

	engine.GetTimingWheel().Start()
	return nil
}

func InitAgent() error {
	svr := agent.NewServer()
	svr.Start()

	return nil
}

func register(param map[string]interface{}) error {
	svrOpt := param["opt"].(*server.Options)
	err := engine.GetServiceContext().GetRegistry().Register(svrOpt)
	if err != nil {
		logger.Error("service registry fail")
		return err
	}

	ttl := param["ttl"].(time.Duration)
	if ttl > 0 {
		engine.GetTimingWheel().NewWheel(map[string]interface{}{"opt": svrOpt, "ttl": ttl}, ttl, register)
	}

	return nil
}