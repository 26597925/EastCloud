package api

import (
	"errors"
	"sapi/pkg/config/encoder"
	"sapi/pkg/config/encoder/hcl"
	"sapi/pkg/config/encoder/json"
	"sapi/pkg/config/encoder/toml"
	"sapi/pkg/config/encoder/yaml"
)

//通过type来设置是否来读取配置，文件就不读线上的，
type Watcher interface {
	Next() (interface{}, error)
	Stop() error
}

var (
	// ErrWatcherStopped is returned when source watcher has been stopped
	ErrWatcherStopped = errors.New("watcher stopped")
	Encoders = map[string]encoder.Encoder{
		"json": json.NewEncoder(),
		"yaml": yaml.NewEncoder(),
		"toml": toml.NewEncoder(),
		"hcl":  hcl.NewEncoder(),
		"yml":  yaml.NewEncoder(),
	}
)