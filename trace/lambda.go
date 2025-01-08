package trace

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go.opentelemetry.io/otel/trace"
)

// StartLambdaSpan extracts data from request, store in the context, and starts a span.
// ContextWithHeaders for propagation of traceContext and baggage;
// ContextWithRequest for api gateway request context.
func StartLambdaSpan(ctx context.Context, request events.APIGatewayProxyRequest) (context.Context, trace.Span) {
	tp := GetTracerProvider()
	tracer := tp.Tracer("")
	logger.Infof("[StartLambdaSpan ctx]: %v\n", ctx)

	ctx = contextWithLambdaRequest(ctx, request)

	spanName := fmt.Sprintf("%s %s", request.HTTPMethod, request.Path)
	ctx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))
	logger.Infof("[StartLambdaSpan tracer.Start ctx]: %v\n", ctx)
	ShowSpanDetails(ctx)
	ShowBaggage(ctx)

	return ctx, span
}

func contextWithLambdaRequest(ctx context.Context, request events.APIGatewayProxyRequest) context.Context {
	ctx = ContextWithHeaders(ctx, request.Headers)
	logger.Debugf("ContextWithHeaders: %v\n", ctx)

	ctx = ContextWithBaggage(ctx, request.Headers)
	logger.Debugf("ContextWithBaggage: %v\n", ctx)

	ctx = ContextWithRequest(ctx, request)
	logger.Debugf("ContextWithRequest: %v\n", ctx)

	return ctx
}
