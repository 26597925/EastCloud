package etcd

import (
	"context"
	"encoding/json"
	"errors"
	retcd "github.com/coreos/etcd/clientv3"
	"sapi/pkg/registry"
)

type etcdWatcher struct {
	stop    chan bool
	w       retcd.WatchChan
	client  *retcd.Client
}

func newEtcdWatcher(r *etcdRegistry) (registry.Watcher, error) {
	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan bool, 1)

	go func() {
		<-stop
		cancel()
	}()

	return &etcdWatcher{
		stop:    stop,
		w:       r.client.Watch(ctx, r.options.Prefix, retcd.WithPrefix(), retcd.WithPrevKV()),
		client:  r.client,
	}, nil
}

func (ew *etcdWatcher) Next() (*registry.Result, error) {
	for res := range ew.w {
		if res.Err() != nil {
			return nil, res.Err()
		}
		if res.Canceled {
			return nil, errors.New("could not get next")
		}
		for _, ev := range res.Events {
			var service *registry.Service
			json.Unmarshal(ev.Kv.Value, &service)
			var action string

			switch ev.Type {
			case retcd.EventTypePut:
				if ev.IsCreate() {
					action = "create"
				} else if ev.IsModify() {
					action = "update"
				}
			case retcd.EventTypeDelete:
				action = "delete"

				json.Unmarshal(ev.PrevKv.Value, &service)
			}

			if service == nil {
				continue
			}
			return &registry.Result{
				Action:  action,
				Service: service,
			}, nil
		}
	}
	return nil, errors.New("could not get next")
}

func (ew *etcdWatcher) Stop() {
	select {
	case <-ew.stop:
		return
	default:
		close(ew.stop)
	}
}