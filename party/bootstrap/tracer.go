package bootstrap

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// NewTracerProvider ...
func NewTracerProvider(serviceInfo *ServiceInfo) (*sdktrace.TracerProvider, error) {
	//traceUrl := strings.ReplaceAll(viper.GetString(env.PathTraceEndpoint), "{{env}}", env.GetMode())
	//exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(traceUrl)))
	//if err != nil {
	//	return nil, err
	//}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1.0))), // OPTIMIZE: Ratio
		//traceSdk.WithBatcher(exp),
		sdktrace.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(serviceInfo.Name),
			semconv.ServiceVersionKey.String(serviceInfo.Version),
			semconv.ServiceInstanceIDKey.String(serviceInfo.Id),
			// attribute.String("env", env.GetMode()),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp, nil
}

//// NewAliyunSlsTracerProvider ...
//func NewAliyunSlsTracerProvider(serviceInfo *ServiceInfo) (*traceSdk.TracerProvider, error) {
//	traceUrl := viper.GetString(env.PathSlsTraceUrl)
//	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(traceUrl)))
//	if err != nil {
//		return nil, err
//	}
//	tp := traceSdk.NewTracerProvider(
//		traceSdk.WithSampler(traceSdk.ParentBased(traceSdk.TraceIDRatioBased(1.0))), // OPTIMIZE: Ratio
//		traceSdk.WithBatcher(exp),
//		traceSdk.WithResource(
//			resource.NewSchemaless(
//				semConv.ServiceNameKey.String(serviceInfo.Name),
//				semConv.ServiceVersionKey.String(serviceInfo.Version),
//				semConv.ServiceInstanceIDKey.String(serviceInfo.Id),
//				attribute.String("sls.otel.project", viper.GetString(env.PathSlsTraceProject)),
//				attribute.String("sls.otel.instanceid", viper.GetString(env.PathSlsTraceInstanceId)),
//				attribute.String("sls.otel.akid", env.GetAliyunAccessKey()),
//				attribute.String("sls.otel.aksecret", env.GetAliyunSecretKey()),
//				// attribute.String("env", env),
//			),
//		),
//	)
//	otel.SetTracerProvider(tp)
//	return tp, nil
//}
