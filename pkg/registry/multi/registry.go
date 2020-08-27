package multi

import (
	"golang.org/x/sync/errgroup"
	"sapi/pkg/registry"
	"sapi/pkg/server/api"
)

type multiRegistry struct {
	registries []registry.Registry
}

func New(registries ...registry.Registry) registry.Registry {
	return &multiRegistry{
		registries: registries,
	}
}

func (m *multiRegistry) Register(option api.Option) error {
	var eg errgroup.Group
	for _, registry := range m.registries {
		eg.Go(func() error {
			return registry.Register(option)
		})
	}
	return eg.Wait()
}

func (m *multiRegistry) Deregister(sv *registry.Service) error {
	var eg errgroup.Group
	for _, registry := range m.registries {
		eg.Go(func() error {
			return registry.Deregister(sv)
		})
	}
	return eg.Wait()
}

func (m *multiRegistry) GetService(name string) (sv *registry.Service, err error) {
	var eg errgroup.Group
	for _, registry := range m.registries {
		eg.Go(func() error {
			sv, err = registry.GetService(name)
			return err
		})
	}
	return sv, eg.Wait()
}

func (m *multiRegistry)  ListServices() (sv []*registry.Service, err error) {
	var eg errgroup.Group
	for _, registry := range m.registries {
		eg.Go(func() error {
			sv, err = registry.ListServices()
			return err
		})
	}
	return sv, eg.Wait()
}

func (m *multiRegistry)  Watch() (registry.Watcher, error) {
	var watchers []registry.Watcher
	var eg errgroup.Group
	for _, registry := range m.registries {
		eg.Go(func() error {
			watch, err := registry.Watch()
			if err == nil {
				watchers = append(watchers, watch)
			}
			return err
		})
	}

	return NewWatcher(watchers...), eg.Wait()
}

func (m *multiRegistry) Close() error {
	var eg errgroup.Group
	for _, registry := range m.registries {
		eg.Go(func() error {
			return registry.Close()
		})
	}
	return eg.Wait()
}

func (m *multiRegistry)  Type() string {
	return "multi"
}