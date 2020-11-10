package multi

import (
	"github.com/26597925/EastCloud/pkg/registry"
	"golang.org/x/sync/errgroup"
)

type multiWatcher struct {
	watchers []registry.Watcher
}

func NewWatcher(watchers ...registry.Watcher) registry.Watcher {
	return &multiWatcher{
		watchers: watchers,
	}
}

func (mw *multiWatcher) Next() (res *registry.Result, err error) {
	var eg errgroup.Group
	for _, watcher := range mw.watchers {
		eg.Go(func() error {
			res, err = watcher.Next()
			return nil
		})
	}
	return res, eg.Wait()
}

func (mw *multiWatcher) Stop() {
	var eg errgroup.Group
	for _, watcher := range mw.watchers {
		eg.Go(func() error {
			watcher.Stop()
			return nil
		})
	}
	eg.Wait()
}
