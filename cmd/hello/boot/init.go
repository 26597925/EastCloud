package boot

import (
	"sapi/cmd/hello/boot/engine"
	"sapi/cmd/hello/router"
	"sapi/pkg/bootstrap"
	"sapi/pkg/bootstrap/flag"
	"sapi/pkg/client/etcdv3"
	"sapi/pkg/client/redis"
	capi "sapi/pkg/config/api"
	"sapi/pkg/logger"
	"sapi/pkg/model"
	"sapi/pkg/registry"
	"sapi/pkg/registry/etcd"
	"sapi/pkg/registry/multi"
	"sapi/pkg/server"
	"sapi/pkg/server/api"
	"sapi/pkg/tracer"
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
	cli := etcdv3.NewOptions().Build()
	opt := &registry.Options{Timeout:3, TTL: 5}
	rsy, err := etcd.NewRegistry(opt, cli)
	if err != nil {
		return err
	}

	engine.SetRegistry(multi.New(rsy))
	for _, svrOpt := range engine.GetConfig().Server {
		err = engine.GetServiceContext().GetRegistry().Register(svrOpt)
		if err != nil {
			return err
		}
	}
	return nil
}