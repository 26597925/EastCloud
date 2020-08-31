package registry

import "sapi/pkg/server/api"

var (
	DefaultPrefix = "/sapi/registry"
)

type Registry interface {
	Register(option api.Option) error
	Deregister(sv *Service) error
	GetService(name string) (*Service, error)
	ListServices() ([]*Service, error)
	Watch() (Watcher, error)
	Close() error
	Type() string
}

type Service struct {
	Driver    string
	Name      string
	ID		  string
	Version   string
	Region    string
	Zone      string
	GroupName string
	IP        string
	Port      int
}

