package hook

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
	"strings"
)

type RedisHook struct{}

func (RedisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx, nil
	}

	tracer := global.Tracer("github.com/go-redis/redis")
	ctx, span := tracer.Start(ctx, "Redis::" + cmd.FullName())
	span.SetAttributes(
		label.String("redis.cmd", cmd.String()),
	)

	return ctx, nil
}

func (RedisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	trace.SpanFromContext(ctx).End()
	return nil
}

func (RedisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx, nil
	}

	names := make([]string, 0, len(cmds))
	data := make([]string, 0, len(cmds))
	for _, cmd := range cmds {
		names = append(names, "Redis::" + cmd.FullName())
		data = append(data, cmd.String())
	}
	tracer := global.Tracer("github.com/go-redis/redis")
	ctx, span := tracer.Start(ctx, strings.Join(names, ","))
	span.SetAttributes(
		label.String("redis", "pipeline"),
		label.Int("redis.num_cmd", len(cmds)),
		label.String("redis.cmds", strings.Join(data, ", ")),
	)

	return ctx, nil
}

func (RedisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	trace.SpanFromContext(ctx).End()
	return nil
}