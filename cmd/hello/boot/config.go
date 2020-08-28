package boot

import (
	"github.com/imdario/mergo"
	flags "sapi/pkg/bootstrap/flag"
	"sapi/pkg/client/redis"
	"sapi/pkg/config"
	"sapi/pkg/config/api"
	"sapi/pkg/config/env"
	"sapi/pkg/config/etcd"
	"sapi/pkg/config/file"
	"sapi/pkg/config/flag"
	"sapi/pkg/logger"
	"sapi/pkg/model"
	"sapi/pkg/registry"
	"sapi/pkg/server"
	"sapi/pkg/tracer"
	"sapi/pkg/util/fileext"
)

type Config struct {
	Name string
	Mode string

	Config *config.Options
	Logger *logger.Options
	Redis *redis.Options
	Orm *model.Options
	Tracer *tracer.Options
	Server []*server.Options
	Registry *registry.Options

	Watcher api.Watcher
}


func parseConfig(fs *flags.Set) (*Config, error) {
	path, err := parsePath(fs)
	if err != nil {
		return nil, err
	}

	var config Config
	err = file.Parse(path, &config)
	if err != nil {
		return nil, err
	}

	var flagConfig Config
	if len(config.Config.FlagPrefixes) > 0 {
		var b []byte
		b, err = flag.Parse(fs, config.Config.FlagPrefixes)
		api.Encoders["json"].Decode(b, &flagConfig)
	}
	mergo.Merge(&flagConfig, config)

	var envConfig Config
	if len(config.Config.EnvPrefixes) > 0 {
		var b []byte
		b, err = env.Parse(config.Config.EnvPrefixes)
		api.Encoders["json"].Decode(b, &envConfig)
	}
	mergo.Merge(&envConfig, flagConfig)

	var etcdConfig Config
	var ed *etcd.Etcd
	if envConfig.Config.Type == "online" {
		ed, err = etcd.NewEtcd(envConfig.Config.Etcd)
		if err != nil {
			return nil, err
		}

		err = ed.Read(&etcdConfig)
		if err != nil {
			return nil , err
		}
	}
	mergo.Merge(&etcdConfig, envConfig)

	var watch api.Watcher
	if etcdConfig.Config.Watcher && etcdConfig.Config.Type == "local" {
		watch, err = file.Watch(path)
	}

	if envConfig.Config.Watcher && envConfig.Config.Type == "online"{
		watch, err = ed.Watch()
	}

	if err != nil {
		return nil, err
	}
	etcdConfig.Watcher = watch

	return  &etcdConfig, nil
}

func parsePath(fs *flags.Set) (string, error) {
	path := fs.String("config")
	exists, err:= fileext.PathExists(path)

	if !exists {
		paths := []string {
			"config/config.yml",
			"config/config.json",
			"config/config.hcl",
			"config/config.toml",
		}
		for _, ph := range paths {
			exists, _ = fileext.PathExists(ph)
			if exists {
				path = ph
				break
			}
		}
	}

	if path == "" {
		return path, err
	}

	return path, nil
}