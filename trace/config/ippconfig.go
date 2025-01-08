package mppconfig

import (
	"context"
	"fmt"

	servicedetector "github.com/aloong-planet/opentelemetry-lambda/trace/detector"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/otel/sdk/resource"

	lambdadetector "go.opentelemetry.io/contrib/detectors/aws/lambda"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// DetectResource returns a resource.Resource merged with two parts:
// 1. configured with the lambda resource detector
// 2. configured with the customized resource detector
func DetectResource(ctx context.Context) (*resource.Resource, error) {
	fmt.Println("detectResource")
	// detect the lambda resource (faas)
	lambdaDetector := lambdadetector.NewResourceDetector()
	lambdaResource, err := lambdaDetector.Detect(ctx)
	if err != nil {
		return nil, err
	}

	// detect the lambda function's environment variables
	serviceDetector := servicedetector.NewServiceDetector()
	svcResource, err := serviceDetector.Detect(ctx)
	if err != nil {
		return nil, err
	}

	// merge the two resources
	res, err := resource.Merge(lambdaResource, svcResource)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func headerCarrier([]byte) propagation.TextMapCarrier {
	fmt.Println("x-trace-id headerCarrier called")
	return propagation.HeaderCarrier{"x-trace-id": []string{"11111"}}
}

func WithHeaderCarrier() otellambda.Option {
	fmt.Println(" WithHeaderCarrier called")
	return otellambda.WithEventToCarrier(headerCarrier)
}

func WithTraceContextPropagator() otellambda.Option {
	return otellambda.WithPropagator(propagation.TraceContext{})
}

// WithDefaultOptions returns a list of all otellambda.Option(s)
// for the otellambda package when using OpenTelemetry.
func WithDefaultOptions(tp *sdktrace.TracerProvider) []otellambda.Option {
	//return []otellambda.Option{WithHeaderCarrier(), WithTraceContextPropagator(), otellambda.WithTracerProvider(tp), otellambda.WithFlusher(tp)}
	return []otellambda.Option{otellambda.WithTracerProvider(tp), otellambda.WithFlusher(tp)}
}
