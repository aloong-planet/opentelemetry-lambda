package trace

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
)

type requestKey struct{}

// Declare a variable of type requestKey. This variable will be used as the key to store and retrieve values
// from the context.
var reqKey = requestKey{}

// RequestContext is a struct that contains the requestContext from the APIGatewayProxyRequest.
type RequestContext struct {
	GatewayProxyContext events.APIGatewayProxyRequestContext
	Headers             map[string]string
}

// ContextWithRequest creates a new context with the requestContext stored in it.
// ctx: WithValue(type trace.requestKey, val <not Stringer>)
func ContextWithRequest(ctx context.Context, apiGwRequest events.APIGatewayProxyRequest) context.Context {
	rc := RequestContext{GatewayProxyContext: apiGwRequest.RequestContext, Headers: apiGwRequest.Headers}
	logger.Debugf("[ContextWithRequest] rc: %+v\n", rc)
	ctx = context.WithValue(ctx, reqKey, rc)

	return ctx
}

// GetRequestContextFromContext retrieves the requestContext from the context.
func GetRequestContextFromContext(ctx context.Context) (RequestContext, bool) {
	rc, ok := ctx.Value(reqKey).(RequestContext)
	if !ok {
		logger.Info("[GetRequestContextFromContext] requestContext not found in context")
	}
	logger.Debugf("[GetRequestContextFromContext] requestContext: %+v\n", rc)

	return rc, ok
}
