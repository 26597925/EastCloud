package main

import (
	"github.com/26597925/EastCloud/cmd/hello/boot"
	"github.com/26597925/EastCloud/pkg/bootstrap"
	"log"

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
