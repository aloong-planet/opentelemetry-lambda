package trace

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

type RequestSpanProcessor struct {
	// keep a reference to the next processor in the chain
	next trace.SpanProcessor
}

// OnStart Process the spans and add attributes from the requestContext to the span.
// The requestContext stores contexts of events.APIGatewayProxyRequestContext.
func (p *RequestSpanProcessor) OnStart(parent context.Context, s trace.ReadWriteSpan) {
	if rc, ok := GetRequestContextFromContext(parent); ok {
		// attributes from ReqContexts
		s.SetAttributes(
			attribute.String("http.method", rc.GatewayProxyContext.HTTPMethod),
			attribute.String("http.path", rc.GatewayProxyContext.Path),
			attribute.String("http.sourceIP", rc.GatewayProxyContext.Identity.SourceIP),
			attribute.String("http.user_agent", rc.GatewayProxyContext.Identity.UserAgent),
			attribute.String("http.accountID", rc.GatewayProxyContext.AccountID),
			attribute.String("http.api_id", rc.GatewayProxyContext.APIID),
			attribute.String("http.api_stage", rc.GatewayProxyContext.Stage),
			attribute.String("http.requestID", rc.GatewayProxyContext.RequestID),
			attribute.String("http.requestTime", rc.GatewayProxyContext.RequestTime),
			attribute.String("http.resourceID", rc.GatewayProxyContext.ResourceID),
			attribute.String("http.resourcePath", rc.GatewayProxyContext.ResourcePath),
			attribute.String("http.protocol", rc.GatewayProxyContext.Protocol),
		)
		// attributes from Headers
		s.SetAttributes(
			attribute.String("x-trace-id", rc.Headers["x-trace-id"]),
		)
	} else {
		logger.Warn("no requestContext found in context")
	}

	// attributes from baggage
	attrs := AttributesFromBaggage(parent)
	s.SetAttributes(attrs...)

	// Pass the event to the next processor
	if p.next != nil {
		p.next.OnStart(parent, s)
	}
}

func (p *RequestSpanProcessor) OnEnd(e trace.ReadOnlySpan) {
	if p.next != nil {
		p.next.OnEnd(e)
	}
}

func (p *RequestSpanProcessor) Shutdown(ctx context.Context) error {
	if p.next != nil {
		return p.next.Shutdown(ctx)
	}
	return nil
}

func (p *RequestSpanProcessor) ForceFlush(ctx context.Context) error {
	if p.next != nil {
		return p.next.ForceFlush(ctx)
	}
	return nil
}
