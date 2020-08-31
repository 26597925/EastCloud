package engine

import (
	"context"
	"sapi/pkg/bootstrap/flag"
	"sapi/pkg/registry"
	"sapi/pkg/server/api"
)

var (
	sc = newServiceContext()
)

type ServiceContext struct {
	ctx context.Context
	cf *Config

	fg []flag.Flag
	fs *flag.Set

	srv []api.Server
	rgy registry.Registry
}

func newServiceContext() *ServiceContext {
	return &ServiceContext{
		ctx: context.Background(),
		srv: make([]api.Server, 0),
	}
}

func GetContext() context.Context {
	return sc.ctx
}

func GetServiceContext() *ServiceContext{
	return sc
}

func AddServe(s ...api.Server) {
	sc.srv = append(sc.srv, s...)
}

func SetFlag(fg []flag.Flag) {
	sc.fg = fg
}

func GetFlag() []flag.Flag {
	return sc.fg
}

func SetFlagSet(fs *flag.Set) {
	sc.fs = fs
}

func GetFlagSet() *flag.Set{
	return sc.fs
}

func SetConfig(cf *Config) {
	sc.cf = cf
}

func GetConfig() *Config {
	return sc.cf
}

func SetRegistry(rgy registry.Registry) {
	sc.rgy = rgy
}

func (sc *ServiceContext) GetServers() []api.Server {
	return sc.srv
}

func (sc *ServiceContext) GetRegistry() registry.Registry {
	return sc.rgy
}