package engine

import (
	"context"
	"sapi/pkg/bootstrap/flag"
	"sapi/pkg/registry"
	"sapi/pkg/server/api"
	"sapi/pkg/util/timer"
	"time"
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

	tw *timer.TimingWheel
}

func newServiceContext() *ServiceContext {
	tw := timer.NewTimingWheel(time.Millisecond * 10)
	return &ServiceContext{
		ctx: context.Background(),
		srv: make([]api.Server, 0),
		tw: tw,
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

func GetTimingWheel() *timer.TimingWheel {
	return sc.tw
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