package bootstrap

import (
	"context"
	"github.com/26597925/EastCloud/pkg/bootstrap/signal"
	"github.com/26597925/EastCloud/pkg/logger"
	"github.com/26597925/EastCloud/pkg/registry"
	"github.com/26597925/EastCloud/pkg/server/api"
	"github.com/26597925/EastCloud/pkg/util/goext"
	"golang.org/x/sync/errgroup"
	"sync"
)

type EngineContext interface {
	GetServers() []api.Server
	GetRegistry() registry.Registry
}

type Engine struct {
	registry registry.Registry

	ec        EngineContext
	initOnce  sync.Once
	startOnce sync.Once
	stopOnce  sync.Once

	signalHooker func(app signal.App)
}

func NewEngine(ec EngineContext) *Engine{
	eng := &Engine{
		ec: ec,
	}
	eng.initialize()
	return eng
}

func (eng *Engine) initialize() {
	eng.initOnce.Do(func() {
		eng.signalHooker = signal.HookSignals
	})

}

func (eng *Engine) Startup (fns ...func() error) (err error){
	eng.startOnce.Do(func() {
		eng.signalHooker(eng)
		err = goext.SerialUntilError(fns...)()
	})

	if err != nil {
		return err
	}

	return nil
}

func (eng *Engine) SetRegistry(reg registry.Registry) {
	eng.registry = reg
}

func (eng *Engine) Serve () (err error) {
	var eg errgroup.Group
	for _, svr := range eng.ec.GetServers() {
		s := svr
		eg.Go(func() error {
			return s.Start()
		})
	}

	logger.Info("start successfully")
	return eg.Wait()
}

func (eng *Engine) Stop () (err error){
	eng.stopOnce.Do(func() {
		var eg errgroup.Group
		for _, svr := range eng.ec.GetServers() {
			s := svr
			eg.Go(s.Stop)
		}
		err = eg.Wait()
	})
	return
}

func (eng *Engine) GracefulStop(ctx context.Context) (err error) {
	eng.stopOnce.Do(func() {
		var eg errgroup.Group
		for _, svr := range eng.ec.GetServers() {
			s := svr
			eg.Go(func() error {
				return s.GracefulStop(ctx)
			})
		}
		err = eg.Wait()
	})
	return err
}

func (eng *Engine) SetSignalHooker(hook func(signal.App)) {
	eng.signalHooker = hook
}