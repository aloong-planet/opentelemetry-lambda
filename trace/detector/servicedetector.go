package servicedetector

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// For a complete list of reserved environment variables in Lambda, see:
// https://docs.aws.amazon.com/lambda/latest/dg/configuration-envvars.html
const (
	stageEnvVar              = "STAGE"
	serviceNameEnvVar        = "MPP_OTEL_SERVICE_NAME"
	lambdaFunctionNameEnvVar = "AWS_LAMBDA_FUNCTION_NAME"
	awsRegionEnvVar          = "AWS_REGION"
)

var (
	empty          = resource.Empty()
	errNotOnLambda = errors.New("process is not on Lambda, cannot detect environment variables from Lambda")
)

// resource detector collects resource information from Lambda environment.
type serviceDetector struct{}

// compile time assertion that resource detector implements the resource.Detector interface.
var _ resource.Detector = (*serviceDetector)(nil)

// NewServiceDetector returns a resource detector that will detect AWS Lambda customized service resources.
func NewServiceDetector() resource.Detector {
	return &serviceDetector{}
}

// Detect collects customized environment variables and generate attributes for mpp.
func (detector *serviceDetector) Detect(context.Context) (*resource.Resource, error) {
	// Lambda resources come from ENV
	lambdaName := os.Getenv(lambdaFunctionNameEnvVar)
	if len(lambdaName) == 0 {
		return empty, errNotOnLambda
	}

	// customize service name
	stage := os.Getenv(stageEnvVar)
	serviceNameIdentifier := os.Getenv(serviceNameEnvVar)
	awsRegion := os.Getenv(awsRegionEnvVar)
	serviceName := fmt.Sprintf("%s_%s_%s_aws_%s", stage, serviceNameIdentifier, lambdaName, awsRegion)
	log.Printf("serviceName: %s\n", serviceName)
	attrs := []attribute.KeyValue{
		semconv.ServiceName(serviceName),
	}

	return resource.NewWithAttributes(semconv.SchemaURL, attrs...), nil
}
