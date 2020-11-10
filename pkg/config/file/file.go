package file

import (
	"errors"
	"github.com/26597925/EastCloud/pkg/config/api"
	"github.com/26597925/EastCloud/pkg/util/fileext"
	"io/ioutil"
	"os"
)

func Parse(path string, data interface{}) error {
	fh, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fh.Close()
	b, err := ioutil.ReadAll(fh)
	if err != nil {
		return err
	}

	if api.Encoders[fileext.Format(path)] != nil {
		err = api.Encoders[fileext.Format(path)].Decode(b, data)
	} else {
		return errors.New("handler errors is nil ")
	}

	return err
}

func Watch(path string) (api.Watcher, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	return newWatcher(path)
}