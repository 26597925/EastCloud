package boot

import (
	"context"
	"sapi/pkg/bootstrap/flag"
	"sapi/pkg/registry"
	"sapi/pkg/server/api"
)

var (
	sc = NewServiceContext()
)

type ServiceContext struct {
	ctx context.Context
	cf *Config

	fg []flag.Flag
	fs *flag.Set

	srv []api.Server
	rgy registry.Registry
}

func NewServiceContext() *ServiceContext {
	return &ServiceContext{
		ctx: context.Background(),
		srv: make([]api.Server, 0),
	}
}

func (sc *ServiceContext) GetServers() []api.Server {
	return sc.srv
}

func (sc *ServiceContext) GetRegistry() registry.Registry {
	return sc.rgy
}

func AddServe(s ...api.Server) {
	sc.srv = append(sc.srv, s...)
}

func GetContext() context.Context {
	return sc.ctx
}

func GetConfig() *Config {
	return sc.cf
}

func Registry()  registry.Registry{
	return sc.rgy
}