package redis

import (
	"context"
	"fmt"

	"gitlab.com/a16624741591/zero/core/env"

	red "github.com/go-redis/redis/v8"
	"gitlab.com/a16624741591/zero/core/trace"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type tracingHook struct {
}

func (t tracingHook) BeforeProcess(ctx context.Context, cmd red.Cmder) (context.Context, error) {
	var hasParentTraceID = oteltrace.SpanContextFromContext(ctx).HasTraceID()
	if hasParentTraceID {
		tracer := trace.GetTracerProvider(trace.RedisTraceName)
		spanCtx, span := tracer.Start(ctx, "[Redis]",
			oteltrace.WithSpanKind(oteltrace.SpanKindClient),
		)
		span.AddEvent("request", oteltrace.WithAttributes(
			attribute.String("req.cmd", cmd.String())))
		return spanCtx, nil
	}
	return ctx, nil
}

func (t tracingHook) AfterProcess(ctx context.Context, cmd red.Cmder) error {
	span := oteltrace.SpanFromContext(ctx)
	if span != nil {
		span.End()
	}
	return nil
}

func (t tracingHook) BeforeProcessPipeline(ctx context.Context, cmds []red.Cmder) (context.Context, error) {
	var hasParentTraceID = oteltrace.SpanContextFromContext(ctx).HasTraceID()
	if hasParentTraceID && env.EnableTracing() {
		tracer := trace.GetTracerProvider(trace.RedisTraceName)
		spanCtx, span := tracer.Start(ctx, "[Redis Pipeline]", oteltrace.WithSpanKind(oteltrace.SpanKindClient))

		var eventOptions = make([]attribute.KeyValue, 0)
		for idx, cmd := range cmds {
			eventOptions = append(eventOptions, attribute.String(fmt.Sprintf("req.cmd%d", idx), cmd.String()))
		}
		span.AddEvent("request", oteltrace.WithAttributes(eventOptions...))
		return spanCtx, nil
	}
	return ctx, nil
}

func (t tracingHook) AfterProcessPipeline(ctx context.Context, cmds []red.Cmder) error {
	span := oteltrace.SpanFromContext(ctx)
	if span != nil {
		span.End()
	}
	return nil
}
