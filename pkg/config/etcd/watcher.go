package etcd

import (
	"context"
	cetcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"sapi/pkg/config/api"
	"sync"
)

type watcher struct {
	name        string
	prefix      string
	format      string

	sync.RWMutex
	ch   chan interface{}
	exit chan bool
}

func newWatcher(etcd *Etcd, wc cetcd.Watcher) (api.Watcher, error) {

	w := &watcher{
		name:        "etcdv3",
		format:		 etcd.format,
		prefix:      etcd.prefix,
		ch:          make(chan interface{}),
		exit:        make(chan bool),
	}

	ch := wc.Watch(context.Background(), w.prefix)

	go w.run(wc, ch)

	return w, nil
}

func (w *watcher) handle(evs []*cetcd.Event) {
	var kv *mvccpb.KeyValue
	for _, v := range evs {
		switch mvccpb.Event_EventType(v.Type) {
		case mvccpb.DELETE:
			return
		default:
			kv = (*mvccpb.KeyValue)(v.Kv)
		}
	}

	var data interface{}
	err := api.Encoders[w.format].Decode(kv.Value, &data)

	if err != nil {
		return
	}

	w.ch <- data
}

func (w *watcher) run(wc cetcd.Watcher, ch cetcd.WatchChan) {
	for {
		select {
		case rsp, ok := <-ch:
			if !ok {
				return
			}
			w.handle(rsp.Events)
		case <-w.exit:
			wc.Close()
			return
		}
	}
}

func (w *watcher) Next() (interface{}, error) {
	select {
	case cs := <-w.ch:
		return cs, nil
	case <-w.exit:
		return nil, api.ErrWatcherStopped
	}
}

func (w *watcher) Stop() error {
	select {
	case <-w.exit:
		return nil
	default:
		close(w.exit)
	}
	return nil
}
