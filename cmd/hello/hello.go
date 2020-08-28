package main

import (
	"log"
	"sapi/cmd/hello/boot"
	"sapi/pkg/bootstrap"
)

func main() {
	eng := bootstrap.NewEngine(boot.Init())
	err := eng.Startup(
		boot.InitFlag,
		boot.InitConfig,
		boot.InitLog,
		boot.InitRedis,
		boot.InitModel,
		boot.InitTracer,
		boot.InitServer,
		boot.InitRegistry,
		)

	if err != nil {
		log.Panic(err)
		return
	}

	eng.Serve()

}
