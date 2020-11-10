package file

import (
	"github.com/26597925/EastCloud/pkg/config/api"
	"github.com/fsnotify/fsnotify"
	"os"
	"runtime"
	"sync"
)

type watcher struct {
	sync.RWMutex

	path string

	fw   *fsnotify.Watcher
	exit chan bool
}

func newWatcher(path string) (api.Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fw.Add(path)

	return &watcher{
		path: path,
		fw:   fw,
		exit: make(chan bool),
	}, nil
}

func (w *watcher) Next() (interface{}, error) {
	// is it closed?
	select {
	case <-w.exit:
		return nil, api.ErrWatcherStopped
	default:
	}

	// try get the event
	select {
	case event, _ := <-w.fw.Events:
		if event.Op == fsnotify.Rename {
			// check existence of file, and add watch again
			_, err := os.Stat(event.Name)
			if err == nil || os.IsExist(err) {
				w.fw.Add(event.Name)
			}
		}

		var data interface{}
		w.Lock()
		err := Parse(w.path, &data)
		if err != nil {
			return nil, err
		}
		w.Unlock()

		// add path again for the event bug of fsnotify
		if runtime.GOOS == "linux" {
			w.fw.Add(w.path)
		}

		return data, nil
	case err := <-w.fw.Errors:
		return nil, err
	case <-w.exit:
		return nil, api.ErrWatcherStopped
	}
}

func (w *watcher) Stop() error {
	return w.fw.Close()
}
