package etcd

import (
	"encoding/json"
	"errors"
	"github.com/26597925/EastCloud/pkg/client/etcdv3"
	"github.com/26597925/EastCloud/pkg/registry"
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
		if ev.Type != etcdv3.KeyDelete {
			json.Unmarshal(ev.Value, &service)
		}

		return &registry.Result{
			Type:  ev.Type,
			Key:   ev.Key,
			Service: service,
		}, nil
	}

	return nil, errors.New("could not get next")
}

func (ew *etcdWatcher) Stop() {
	ew.watcher.Stop()
}