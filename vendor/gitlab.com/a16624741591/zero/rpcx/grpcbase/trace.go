package grpcbase

import (
	"context"
	"fmt"

	"github.com/google/martian/log"

	"go.opentelemetry.io/otel/attribute"

	"gitlab.com/a16624741591/zero/core/sysx"
	"gitlab.com/a16624741591/zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func NewTraceSpanFromRequestID(serverTag, rpcServiceName, requestID string) (context.Context, oteltrace.Span, string) {
	if requestID == "" {
		return context.Background(), oteltrace.SpanFromContext(context.Background()), ""
	}
	propagator := otel.GetTextMapPropagator()
	ctx := propagator.Extract(context.Background(), propagation.MapCarrier{trace.OlteTraceHeader: requestID})
	tracer := trace.GetTracerProvider(trace.RpcTraceName)

	name, attr := trace.SpanInfo(rpcServiceName, "")
	spanContext, span := tracer.Start(
		ctx,
		fmt.Sprintf("%s（%s）-> %s", sysx.SubSystem, serverTag, name),
		oteltrace.WithSpanKind(oteltrace.SpanKindClient))
	span.SetAttributes(attr...)
	trace.PatchSpanData(span)

	var mapCarrier = propagation.MapCarrier{}
	propagator.Inject(spanContext, mapCarrier)
	return spanContext, span, mapCarrier[trace.OlteTraceHeader]
}

func RPCServerSideLogic(rpcServiceName, requestID string, requestData interface{}, businessFunc func(ctx context.Context, req interface{}) Response) Response {
	if requestID == "" {
		resp := businessFunc(context.Background(), requestData)
		return resp
	}
	log.Infof("GetUser begin 55555")
	ctx1, span, requestID := NewTraceSpanFromRequestID("server", rpcServiceName, requestID)
	defer span.End()

	//b, _ := jsoniter.Marshal(requestData)
	//span.AddEvent("request", oteltrace.WithAttributes(attribute.String("message.data", string(b))))

	resp := businessFunc(ctx1, requestData)
	span.AddEvent("response", oteltrace.WithAttributes(attribute.String("message.data", resp.String())))
	return resp
}
