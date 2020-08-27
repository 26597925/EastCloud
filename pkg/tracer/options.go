package tracer

type Options struct {
	Driver  string
	Name string
	Zipkin Zipkin
	Jaeger Jaeger
}

type Zipkin struct {
	EndpointUrl string
}

type Jaeger struct {
	Mode string //local,server
	EndpointUrl string
	UserName string
	Password string
	AgentEndpoint string
}

func initOptions() Options{
	option := Options{
		Driver:"jaeger",
		Name:"demo",
		Zipkin: Zipkin {
			EndpointUrl:"http://192.168.16.218:9411/api/v2/spans",
		},
		Jaeger: Jaeger {
			Mode: "server",
			EndpointUrl:"http://192.168.16.233:14268/api/traces",
			AgentEndpoint:"localhost:6831",
		},
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