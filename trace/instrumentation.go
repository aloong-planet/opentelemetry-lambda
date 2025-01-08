package trace

import (
	"github.com/aws/smithy-go/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

// NewHTTPClient returns a new HTTP client instrumented with OpenTelemetry.
func NewHTTPClient() *http.Client {
	return &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}

// AppendAWSOTelMiddlewares wraps the OTel AppendMiddlewares function for the AWS SDK with default configuration.
func AppendAWSOTelMiddlewares(apiOptions *[]func(*middleware.Stack) error) {
	otelaws.AppendMiddlewares(apiOptions)
}
