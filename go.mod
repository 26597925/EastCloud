module sapi

go 1.14

require github.com/gin-gonic/gin v1.5.0

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Shopify/sarama v1.26.4
	github.com/benbjohnson/clock v1.0.3 // indirect
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/coreos/etcd v3.3.22+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fsnotify/fsnotify v1.4.7
	github.com/ghodss/yaml v1.0.0
	github.com/go-redis/redis/v8 v8.0.0-beta.5
	github.com/go-session/session v3.1.2+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/hashicorp/hcl v1.0.0
	github.com/imdario/mergo v0.3.9
	github.com/jinzhu/gorm v1.9.14
	github.com/klauspost/compress v1.10.10 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.13.0
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/tidwall/buntdb v1.1.2
	github.com/tidwall/gjson v1.6.0 // indirect
	github.com/tidwall/pretty v1.0.1 // indirect
	github.com/uber/jaeger-client-go v2.24.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	//go.etcd.io/etcd v2.3.8+incompatible // indirect
	go.opentelemetry.io/otel v0.11.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.11.0
	go.opentelemetry.io/otel/exporters/stdout v0.11.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.11.0
	go.opentelemetry.io/otel/exporters/trace/zipkin v0.11.0
	go.opentelemetry.io/otel/sdk v0.11.0
	go.uber.org/multierr v1.5.0
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/sys v0.0.0-20200824131525-c12d262b63d8 // indirect
	google.golang.org/appengine v1.6.6
	google.golang.org/genproto v0.0.0-20200825200019-8632dd797987 // indirect
	google.golang.org/grpc v1.31.1
	google.golang.org/protobuf v1.25.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace google.golang.org/grpc v1.31.1 => google.golang.org/grpc v1.26.0