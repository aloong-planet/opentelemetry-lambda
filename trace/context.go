package trace

import (
	"context"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

// ContextWithHeaders extracts tracing metadata from incoming request's headers and stores them in the context.
// This is used to propagate the traceContext and baggage from the incoming request to the downstream requests.
// propagation.TraceContext{} header format:
//
//	traceparent: 00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01
//	tracestate:  vendorname=opaquevalue
//
// propagation.Baggage{} header format:
//
//	baggage: key1=val1,key2=val2
//
// ctx: WithValue(type trace.traceContextKeyType, val <not Stringer>)
func ContextWithHeaders(ctx context.Context, headers map[string]string) context.Context {
	h := make(http.Header)
	for k, v := range headers {
		//log.Debug(message)("k: %s, v: %s\n", k, v)
		h.Add(k, v)
	}
	return otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(h))
}

// ContextWithBaggage creates a new baggage context storing customized metadata.
// If the incoming request does not contain the x-trace-id header, treat the current function as the root function.
// create a new 'x-trace-id' and store it in the baggage.
// If the incoming request contains the x-trace-id header, store it in the baggage directly.
// ctx: WithValue(type baggage.baggageContextKeyType, val <not Stringer>)
func ContextWithBaggage(ctx context.Context, headers map[string]string) context.Context {
	xTraceId, ok := headers["x-trace-id"]
	if !ok {
		xTraceId = uuid.New().String()
	}
	rb, _ := baggage.NewMember("x-trace-id", xTraceId)
	m, _ := baggage.New(rb)
	ctx = baggage.ContextWithBaggage(ctx, m)

	return ctx
}

// ShowSpanDetails shows the traceID and spanID of the current span in the context.
func ShowSpanDetails(ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID()
	spanID := span.SpanContext().SpanID()
	logger.Debugf("TraceID: %s, SpanID: %s\n", traceID, spanID)
}

// ShowBaggage shows the key-value pairs of the current baggage in the context.
func ShowBaggage(ctx context.Context) {
	bag := baggage.FromContext(ctx)
	logger.Debug("Baggage: ")
	for _, m := range bag.Members() {
		logger.Debugf("  %s: %s\t", m.Key(), m.Value())
	}
}

// AttributesFromBaggage returns a list of attributes from the current baggage in the context.
func AttributesFromBaggage(ctx context.Context) []attribute.KeyValue {
	bag := baggage.FromContext(ctx)
	attrs := make([]attribute.KeyValue, 0)

	if bag.Len() == 0 {
		return attrs
	}

	for _, m := range bag.Members() {
		logger.Debugf("  %s: %s\n", m.Key(), m.Value())
		attrs = append(attrs, attribute.String(m.Key(), m.Value()))
	}
	return attrs
}
