package server

import (
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

var (
	serverIdPrefix = "server-id-")

type Options struct {
	Name	  string
	Version   string
	Driver    string
	Id        string
	Region    string
	Zone      string
	GroupName string
	IP        string
	Port      int

	GrpcOptions  []grpc.ServerOption
}

func initOptions() *Options{
	id := serverIdPrefix + uuid.New().String()
	option := &Options{
		Driver: "gin",
		Id: id,
		IP: "",
		Port: 8000,
	}

	return option
}

func NewOptions(opts ...Option) *Options {
	options := initOptions()
	for _, o := range opts {
		o(options)
	}

	return options
}

func Driver(driver string) Option {
	return func(o *Options) {
		o.Driver = driver
	}
}

func (opt *Options) GetName() string {
	return opt.Name
}

func (opt *Options) GetDriver() string {
	return opt.Driver
}

func (opt *Options) GetId() string {
	return opt.Id
}

func (opt *Options) GetVersion() string {
	return opt.Version
}

func (opt *Options) GetRegion() string {
	return opt.Region
}

func (opt *Options) GetZone() string {
	return opt.Zone
}

func (opt *Options) GetGroupName() string {
	return opt.GroupName
}

func (opt *Options) GetIP() string {
	return opt.IP
}

func (opt *Options) GetPort() int {
	return opt.Port
}

func (opt *Options) GetGrpcOptions() []grpc.ServerOption {
	return opt.GrpcOptions
}