package model

type Options struct {
	Driver string
	Host string
	Port int
	Database string
	Username string
	Password string
	Charset string
	Prefix string
	DbFile string
}

func initOptions() Options{
	option := Options{
		Driver: "mysql",
		Host:"127.0.0.1",
		Port:3306,
		Database:"",
		Username:"root",
		Password:"",
		Charset:"utf8",
		Prefix:"t_",
	}

	return option
}

func NewOptions(opts ...Option) Options {
	options := initOptions()
	for _, o := range opts {
		o(&options)
	}

	return options
}

func Driver(driver string) Option {
	return func(o *Options) {
		o.Driver = driver
	}
}

func Host(host string) Option {
	return func(o *Options) {
		o.Host = host
	}
}

func Port(port int) Option {
	return func(o *Options) {
		o.Port = port
	}
}

func Database(database string) Option {
	return func(o *Options) {
		o.Database = database
	}
}

func Username(username string) Option {
	return func(o *Options) {
		o.Username = username
	}
}

func Password(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

func Charset(charset string) Option {
	return func(o *Options) {
		o.Charset = charset
	}
}

func Prefix(prefix string) Option {
	return func(o *Options) {
		o.Prefix = prefix
	}
}