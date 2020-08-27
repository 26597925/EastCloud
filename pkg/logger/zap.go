package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type zapLog struct {
	cfg  zap.Config
	zap  *zap.Logger
	opts *Options
}

func loggerToZapLevel(level Level) zapcore.Level {
	switch level {
	case TraceLevel, DebugLevel:
		return zap.DebugLevel
	case InfoLevel:
		return zap.InfoLevel
	case WarnLevel:
		return zap.WarnLevel
	case ErrorLevel:
		return zap.ErrorLevel
	case FatalLevel:
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func NewZap (opts *Options) (Logger, error) {
	var err error

	zapConfig := zap.NewDevelopmentConfig()
	if !opts.Development {
		zapConfig = zap.NewProductionConfig()
	}
	zapConfig.Level.SetLevel(loggerToZapLevel(opts.Level))
	zapConfig.InitialFields = map[string]interface{}{"serviceName": opts.AppName}

	f := func(file string) zapcore.WriteSyncer {
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   opts.LogFileDir + sp + opts.AppName + "-" + file,
			MaxSize:    opts.MaxSize,
			MaxBackups: opts.MaxBackups,
			MaxAge:     opts.MaxAge,
			Compress:   true,
			LocalTime:  true,
		})
	}
	errWS := f(opts.ErrorFileName)
	warnWS := f(opts.WarnFileName)
	infoWS := f(opts.InfoFileName)
	debugWS := f(opts.DebugFileName)
	fileEncoder := zapcore.NewJSONEncoder(zapConfig.EncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(zapConfig.EncoderConfig)

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > zapcore.WarnLevel && zapcore.WarnLevel - zapConfig.Level.Level() > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel - zapConfig.Level.Level() > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel - zapConfig.Level.Level() > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel - zapConfig.Level.Level() > -1
	})

	debugConsoleWS := zapcore.Lock(os.Stdout) // 控制台标准输出
	errorConsoleWS := zapcore.Lock(os.Stderr)
	cores := []zapcore.Core{
		zapcore.NewCore(fileEncoder, errWS, errPriority),
		zapcore.NewCore(fileEncoder, warnWS, warnPriority),
		zapcore.NewCore(fileEncoder, infoWS, infoPriority),
		zapcore.NewCore(fileEncoder, debugWS, debugPriority),
		zapcore.NewCore(consoleEncoder, errorConsoleWS, errPriority),
		zapcore.NewCore(consoleEncoder, debugConsoleWS, warnPriority),
		zapcore.NewCore(consoleEncoder, debugConsoleWS, infoPriority),
		zapcore.NewCore(consoleEncoder, debugConsoleWS, debugPriority),
	}

	op := zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})

	log, err := zapConfig.Build(op, zap.AddCallerSkip(3))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer log.Sync()

	zl := &zapLog{
		cfg:   zapConfig,
		zap:   log,
		opts:  opts,
	}

	return zl, nil
}

func (z *zapLog) Log(level Level, args ...interface{}) {
	data := make([]zap.Field, 0, len(z.opts.fields))
	for k, v := range z.opts.fields {
		data = append(data, zap.Any(k, v))
	}

	lvl := loggerToZapLevel(level)
	msg := fmt.Sprint(args...)

	switch lvl {
	case zap.DebugLevel:
		z.zap.Debug(msg, data...)
	case zap.InfoLevel:
		z.zap.Info(msg, data...)
	case zap.WarnLevel:
		z.zap.Warn(msg, data...)
	case zap.ErrorLevel:
		z.zap.Error(msg, data...)
	case zap.FatalLevel:
		z.zap.Fatal(msg, data...)
	}
}

func (z *zapLog) Logf(level Level, format string, args ...interface{}){
	data := make([]zap.Field, 0, len(z.opts.fields))
	for k, v := range z.opts.fields {
		data = append(data, zap.Any(k, v))
	}

	lvl := loggerToZapLevel(level)
	msg := fmt.Sprintf(format, args...)
	switch lvl {
	case zap.DebugLevel:
		z.zap.Debug(msg, data...)
	case zap.InfoLevel:
		z.zap.Info(msg, data...)
	case zap.WarnLevel:
		z.zap.Warn(msg, data...)
	case zap.ErrorLevel:
		z.zap.Error(msg, data...)
	case zap.FatalLevel:
		z.zap.Fatal(msg, data...)
	}
}

func (z *zapLog) Info(args ...interface{}) {
	z.Log(InfoLevel, args...)
}

func (z *zapLog) InfoF(format string, args ...interface{}) {
	z.Logf(InfoLevel, format, args...)
}

func (z *zapLog) Debug(args ...interface{}) {
	z.Log(DebugLevel, args...)
}

func (z *zapLog) DebugF(format string, args ...interface{}) {
	z.Logf(DebugLevel, format, args...)
}

func (z *zapLog) Warn(args ...interface{}) {
	z.Log(WarnLevel, args...)
}

func (z *zapLog) WarnF(format string, args ...interface{}) {
	z.Logf(WarnLevel, format, args...)
}

func (z *zapLog) Error(args ...interface{}) {
	z.Log(ErrorLevel, args...)
}

func (z *zapLog) ErrorF(format string, args ...interface{}) {
	z.Logf(ErrorLevel, format, args...)
}

func (z *zapLog) Fatal(args ...interface{}) {
	z.Log(FatalLevel, args...)
}

func (z *zapLog) FatalF(format string, args ...interface{}) {
	z.Logf(FatalLevel, format, args...)
}

func (z *zapLog) Type() string {
	return "zap"
}