package config

import "sapi/pkg/config/etcd"

type Options struct {
	Type    string
	Watcher bool
	FlagPrefixes []string
	EnvPrefixes []string
	Etcd    *etcd.Options
}

type Option func(*Options)

func initOptions() *Options{
	return &Options{}
}

func NewOptions(opts ...Option) *Options {
	options := initOptions()
	for _, o := range opts {
		o(options)
	}

	return options
}