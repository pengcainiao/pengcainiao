package handler

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"gitlab.com/a16624741591/zero/core/env"

	"github.com/gin-gonic/gin"
	"gitlab.com/a16624741591/zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// TracingHandler return a middleware that process the opentelemetry.
func TracingHandler(serviceName string) gin.HandlerFunc {
	propagator := otel.GetTextMapPropagator()
	tracer := trace.GetTracerProvider(trace.HttpTraceName)
	return func(r *gin.Context) {
		if strings.Contains(r.Request.URL.Path, "/stream") || strings.Contains(r.Request.URL.Path, "/user/verify") {
			r.Next()
			return
		}

		injectRequestID(r)

		spanName := fmt.Sprintf("[HTTP] %s %s", r.Request.Method, r.Request.URL.Path)
		spanContext, span := tracer.Start(
			r.Request.Context(),
			spanName,
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(
				env.ReleaseMode+":"+serviceName, spanName, r.Request)...),
		)
		defer func() {
			span.End()
		}()
		trace.PatchSpanData(span)

		var mapCarrier = propagation.MapCarrier{}
		propagator.Inject(spanContext, mapCarrier)
		var (
			requestID   = mapCarrier[trace.OlteTraceHeader]
			requestBody []byte
		)
		userID := r.GetHeader("x-auth-user")
		if userID != "" {
			span.SetAttributes(
				attribute.String("messaging.id", requestID),
				attribute.String("user.id", userID),
				attribute.String("user.client_version", r.GetHeader("x-auth-version")),
				attribute.String("user.platform", r.GetHeader("x-auth-platform")))
		}
		r.Writer.Header().Set(trace.TraceIdKey, requestID)
		r.Request.Header.Set(trace.TraceIdKey, requestID)

		if r.Request.Method == http.MethodPost || r.Request.Method == http.MethodPut || r.Request.Method == http.MethodPatch || r.Request.Method == http.MethodDelete {
			requestBody, _ = ioutil.ReadAll(r.Request.Body)
			r.Request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))
		}
		r.Request = r.Request.WithContext(oteltrace.ContextWithSpan(r.Request.Context(), span))
		r.Request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))
		span.AddEvent("request", oteltrace.WithAttributes(
			attribute.String("request.data", string(requestBody)),
			attribute.String("request.query", r.Request.URL.RawQuery),
		))
		r.Next()
	}
}

func injectRequestID(c *gin.Context) {
	reqID := c.GetHeader(trace.TraceIdKey)
	if reqID != "" {
		ctx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.MapCarrier{trace.OlteTraceHeader: reqID})
		c.Request = c.Request.WithContext(ctx)
	}
}
