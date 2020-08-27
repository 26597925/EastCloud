package logger

import (
	"github.com/imdario/mergo"
	"path/filepath"
	"reflect"
)

var(
	sp = string(filepath.Separator)
)

type Options struct {
	Level		  Level
	Development   bool

	LogFileDir    string
	AppName       string
	ErrorFileName string
	WarnFileName  string
	InfoFileName  string
	DebugFileName string
	MaxSize       int
	MaxBackups    int
	MaxAge        int

	fields map[string]interface{}
}

type Transformer struct {
}

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "errors"
	case FatalLevel:
		return "fatal"
	}
	return ""
}

type Option func(*Options)

func initOptions() *Options{
	logFileDir , _ := filepath.Abs(filepath.Dir(filepath.Join(".")))
	logFileDir += sp + "logs" + sp

	options := &Options{
		Level: 0,
		Development: true,
		LogFileDir: logFileDir,
		AppName: "app",
		ErrorFileName: "errors.logger",
		WarnFileName: "warn.logger",
		InfoFileName: "info.logger",
		DebugFileName: "debug.logger",
		MaxSize: 500,// 500M
		MaxBackups: 3,// 3 backup
		MaxAge: 1,// 1 day
		fields: make(map[string]interface{}),
	}

	//MaxSize: 500, 500M MaxBackups: 3, 3 backup MaxAge: 1, 1 day

	return options
}

func NewOptions(opts ...Option) *Options {
	options := initOptions()
	for _, o := range opts {
		o(options)
	}

	return options
}

func (t Transformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	var o Options
	if typ != reflect.TypeOf(o) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				if !reflect.DeepEqual(src.Interface(), reflect.Zero(src.Type()).Interface()) {
					dst.Set(src)
				}
			}
			return nil
		}
	}
	return nil
}


func Merge(opts *Options) Option {
	return func(o *Options) {
		mergo.Merge(o, opts, mergo.WithTransformers(Transformer{}))
	}
}

func Fields(fields map[string]interface{}) Option {
	return func(o *Options) {
		o.fields = fields
	}
}