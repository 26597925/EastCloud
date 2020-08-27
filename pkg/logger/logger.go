package logger

type Level int8

var (
	log, _= NewZap(NewOptions())
	)

const (
	TraceLevel Level = iota - 2
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type Logger interface {
	Log(level Level, args ...interface{})
	Logf(level Level, format string, args ...interface{})
	Info(args ...interface{})
	InfoF(format string, args ...interface{})
	Debug(args ...interface{})
	DebugF(format string, args ...interface{})
	Warn(args ...interface{})
	WarnF(format string, args ...interface{})
	Error(args ...interface{})
	ErrorF(format string, args ...interface{})
	Fatal(args ...interface{})
	FatalF(format string, args ...interface{})
	Type() string
}

func SetLog(logger Logger) {
	log = logger
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func InfoF(format string, args ...interface{}) {
	log.InfoF(format, args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func DebugF(format string, args ...interface{}) {
	log.DebugF(format, args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func WarnF(format string, args ...interface{}){
	log.WarnF(format, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func ErrorF(format string, args ...interface{}) {
	log.ErrorF(format, args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func FatalF(format string, args ...interface{}) {
	log.FatalF(format, args...)
}