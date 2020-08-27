package etcdv3

import (
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"sync"
)

const (
	KeyCreate = iota
	KeyUpdate
	KeyDelete
)

type Event struct {
	Type  int
	Key   string
	Value []byte
}

type watcher struct {
	sync.RWMutex
	ch   chan *Event
	exit chan bool
}

func (w *watcher) Next() (*Event, error) {
	select {
	case evt := <-w.ch:
		return evt, nil
	case <-w.exit:
		return nil, errors.New("watcher stopped")
	}
}

func (w *watcher) Stop()  {
	select {
	case <-w.exit:
		return
	default:
		close(w.exit)
	}
}

func (w *watcher) Type() string {
	return "etcdv3"
}

func (w *watcher) handle(evt *clientv3.Event) {
	event := &Event{
		Key: string(evt.Kv.Key),
		Value: evt.Kv.Value,
	}

	switch evt.Type {
		case mvccpb.PUT:
			if evt.IsCreate() {
				event.Type = KeyCreate
			} else {
				event.Type = KeyUpdate
			}
		case mvccpb.DELETE:
			event.Type = KeyDelete
	}

	w.ch <- event
}

func (w *watcher) run(wc clientv3.Watcher, ch clientv3.WatchChan) {
	for {
		select {
		case rsp, ok := <-ch:
			if !ok {
				return
			}
			for _, v := range rsp.Events {
				w.handle(v)
			}
		case <-w.exit:
			wc.Close()
			return
		}
	}
}
