package etcd

import (
	"encoding/json"
	"errors"
	"sapi/pkg/client/etcdv3"
	"sapi/pkg/registry"
)

type etcdWatcher struct {
	watcher *etcdv3.Watcher
}

func newEtcdWatcher(r *etcdRegistry) (registry.Watcher, error) {
	wr, err := r.client.NewWatchWithPrefixKey(r.opts.Prefix)
	if err != nil {
		return nil, err
	}

	return &etcdWatcher{
		watcher:wr,
	}, nil
}

func (ew *etcdWatcher) Next() (*registry.Result, error) {
	for {
		ev, err := ew.watcher.Next()
		if err != nil {
			return nil, err
		}

		var service *registry.Service
		json.Unmarshal(ev.Value, &service)

		if service == nil {
			continue
		}

		return &registry.Result{
			Type:  ev.Type,
			Service: service,
		}, nil
	}

	return nil, errors.New("could not get next")
}

func (ew *etcdWatcher) Stop() {
	ew.watcher.Stop()
}