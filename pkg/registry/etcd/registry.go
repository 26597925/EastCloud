package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/26597925/EastCloud/pkg/client/etcdv3"
	"github.com/26597925/EastCloud/pkg/registry"
	"github.com/26597925/EastCloud/pkg/server/api"
	"sort"
	"sync"
)

type etcdRegistry struct {
	opts 	  *registry.Options

	client    *etcdv3.Client
	register  sync.Map
}

func NewRegistry(opts *registry.Options, client *etcdv3.Client) (registry.Registry, error) {
	if opts.Prefix == "" {
		opts.Prefix = registry.DefaultPrefix
	}

	e := &etcdRegistry{
		opts: opts,
		client: client,
	}

	return e, nil
}

func (e *etcdRegistry) Register(opt api.Option) (err error) {
	key := fmt.Sprintf("%s/%s/%s/%s", e.opts.Prefix, e.Type(), opt.GetName(), opt.GetId())
	service := &registry.Service{
		Driver:    opt.GetDriver(),
		Name:      opt.GetName(),
		ID:        opt.GetId(),
		Version:   opt.GetVersion(),
		Region:    opt.GetRegion(),
		Zone:      opt.GetZone(),
		GroupName: opt.GetGroupName(),
		IP:        opt.GetIP(),
		Port:      opt.GetPort(),
	}
	val, err := json.Marshal(service)
	if err != nil {
		return err
	}

	err = e.client.NewRegistry().Register(key, string(val), e.opts.TTL)
	if err != nil {
		return err
	}

	e.register.Store(key, val)
	return nil
}

func (e *etcdRegistry) Deregister(sv *registry.Service) error {
	key := fmt.Sprintf("%s/%s/%s/%s", e.opts.Prefix, e.Type(), sv.Name, sv.ID)
	err := e.client.NewRegistry().Deregister(key)
	if err == nil {
		e.register.Delete(key)
	}
	return err
}

func (e *etcdRegistry) GetService(name string) ([]*registry.Service, error) {
	key := fmt.Sprintf("%s/%s/%s", e.opts.Prefix, e.Type(), name)
	rsp, err := e.client.GetPrefix(context.Background(), key)
	if err != nil {
		return nil, err
	}

	services := make([]*registry.Service, 0, len(rsp))
	for _, n := range rsp {
		var s *registry.Service
		err = json.Unmarshal(n, &s)

		if err != nil {
			continue
		}

		services = append(services, s)
	}

	return services, nil
}

func (e *etcdRegistry) ListServices() ([]*registry.Service, error) {
	rsp, err := e.client.GetPrefix(context.Background(), e.opts.Prefix)
	if err != nil {
		return nil, err
	}

	services := make([]*registry.Service, 0, len(rsp))
	for _, n := range rsp {
		var s *registry.Service
		err = json.Unmarshal(n, &s)

		if err != nil {
			continue
		}

		services = append(services, s)
	}

	sort.Slice(services, func(i, j int) bool { return services[i].Name < services[j].Name })

	return services, nil
}

func (e *etcdRegistry) Watch() (registry.Watcher, error) {
	return newEtcdWatcher(e)
}

func (e *etcdRegistry) Close() (err error) {
	var wg sync.WaitGroup
	e.register.Range(func(k, v interface{}) bool {
		wg.Add(1)
		go func(v interface{}) {
			defer wg.Done()
			var s *registry.Service
			err = json.Unmarshal(v.([]byte), &s)
			if err == nil {
				err = e.Deregister(s)
			}
		}(v)
		return true
	})
	wg.Wait()

	return nil
}

func (e *etcdRegistry) Type() string {
	return "etcdv3"
}