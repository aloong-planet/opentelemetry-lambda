# OpenTelemetry Trace
This library is used to enable open-telemetry tracing on aws lambda functions. With the customized lambda collector layer added, you could send data to the monitoring system.

### Environments
|name|type|default|description|
|----------|----------|----------|----------|
|OTEL_TRACE_ENABLED|bool|false|toggle for otel trace lib, default to disabled if not set.|
|OTEL_SERVICE_NAME|string|nil|otel build in environment, stand for service name of the tracing data.|
|OTEL_LOG_LEVEL|string|info|log level for otel trace lib|


## How To Instrument tracing

### Install the trace library

```
go get |github.com:aloong-planet/opentelemetry-lambda/trace
go get go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda
```

### Update lambda function code

```
package main

import (
	"fmt"

	mpp_trace "github.com:aloong-planet/opentelemetry-lambda/trace"
	mppconfig "github.com:aloong-planet/opentelemetry-lambda/trace/config"

	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
)

func lambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (interface{}, error) {
    // your function logic

	// create span here
	ctx, span := mpp_trace.StartLambdaSpan(ctx, request)
	defer span.End()

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       os.Getenv("_X_AMZN_TRACE_ID"),  // your body
	}, nil
}

func main() {
	// init trace provider and instrument lambdaHandler with configs.
	ctx := context.Background()
	tp, _ := mpp_trace.InitTrace(ctx)
	defer mpp_trace.ShutdownTracerProvider(ctx, tp)

	lambda.Start(otellambda.InstrumentHandler(lambdaHandler, mppconfig.WithDefaultOptions(tp)...))
}
```

### Add a customized lambda open-telemetry layer

Configure the lambda function to add a published customized lambda collector layer.

### Enabling the Trace

By default, the trace is not enabled without configuring the environment "OTEL_TRACE_ENABLED=true".

Which means that:
* **No Exporters Configured**: By default, the TracerProvider does not have any exporters configured. This means it will not send trace data (like spans) to any backend system or endpoint. Without exporters, the collected trace data is essentially discarded.
* **No Connection Attempts**: Since there are no exporters configured by default, the TracerProvider will not attempt to connect to any endpoint (like an OTLP endpoint or Jaeger collector).

By adding the lambda function environment `OTEL_TRACE_ENABLED=true`, the library would start a default connection with the open-telemetry collector extension.

### (Optional) Setting trace loglevel

By default, the trace loglevel is 'info', add an environment `OTEL_LOG_LEVEL=debug` to enable debug logs.
