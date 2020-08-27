package boot

import (
	"sapi/cmd/hello/router"
	"sapi/pkg/bootstrap"
	"sapi/pkg/bootstrap/flag"
	"sapi/pkg/client/redis"
	capi "sapi/pkg/config/api"
	"sapi/pkg/logger"
	"sapi/pkg/model"
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
	sc.fg = append(fg, flags...)

	return sc
}

func InitFlag () error {
	fs := flag.NewFlagSet()
	fs.Register(sc.fg...)
	err := fs.Parse()
	if err == nil {
		sc.fs = fs
	}
	return err
}

func InitConfig() error {
	cf, err := parseConfig(sc.fs)
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
	sc.cf = cf
	return err
}

func InitLog() error {
	opts := logger.NewOptions(logger.Merge(GetConfig().Logger))
	log, err := logger.NewZap(opts)

	if err != nil {
		return err
	}

	logger.SetLog(log)
	return nil
}

func InitRedis() error {
	return redis.Init(GetContext(), GetConfig().Redis)
}

func InitModel() error {
	return model.Init(GetConfig().Orm)
}

func InitTracer() error {
	tracer.AddHookSpanCtx(GetContext())
	err := tracer.Init(GetConfig().Tracer)
	return err
}

func InitServer() error {
	 handlers := map[string]api.Handler{
	 	"grpc": router.GrpcRouter,
	 	"http": router.HttpRouter,
	 }

	for _, svrOpt := range GetConfig().Server {
		svr := server.NewServer(svrOpt)
		err := svr.Init()
		if err != nil {
			return err
		}

		svr.Handler(handlers[svr.GetOption().GetName()])
		AddServe(svr)
	}
	return nil
}

func InitRegistry() error {
	return nil
}