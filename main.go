package main

import (
	"net/http"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"

	"github.com/uber/jaeger-client-go"
	jaegerClientConfig "github.com/uber/jaeger-client-go/config"
)

func main() {

	cfg := jaegerClientConfig.Configuration{
		Sampler: &jaegerClientConfig.SamplerConfig{
			Type:              "const",
			Param:             1,
			SamplingServerURL: "localhost:5775",
		},
		Reporter: &jaegerClientConfig.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}
	tracer, closer, _ := cfg.New(
		"some_service",
		jaegerClientConfig.Logger(jaeger.StdLogger),
	)

	defer closer.Close()
	tracer2(testTracer1(tracer)) // working
	// testTracer()

	http.ListenAndServe(":8080", nil)

}

// Working
func testTracer1(t opentracing.Tracer) (opentracing.Tracer, opentracing.Span) {

	span := t.StartSpan("new_span")
	defer span.Finish()
	span.SetOperationName("span_1")
	span.LogFields(log.String("ds", "asd"))
	span.LogEvent("hello")
	span.SetTag(zipkincore.HTTP_PATH, struct{ name string }{"ad"})

	span.SetBaggageItem("Some_Key", "12345")
	span.SetBaggageItem("Some-other-key", "42")

	return t, span
}

func tracer2(t opentracing.Tracer, span opentracing.Span) error {
	span2 := t.StartSpan("span_2", opentracing.ChildOf(span.Context()))
	defer span2.Finish()
	span2.LogFields(log.String("hello", "span2"))
	time.Sleep(1 * time.Millisecond)
	return nil
}

// func testTracer() interface{} {
// 	tracer, closer := jaeger.NewTracer(
// 		"crossdock",
// 		jaeger.NewConstSampler(true),
// 		jaeger.NewNullReporter(),
// 	)
// 	defer closer.Close()

// 	span := tracer.StartSpan("hi")
// 	span.LogEvent("hello")
// 	span.SetBaggageItem("key", "xyz")
// 	// ctx := opentracing.ContextWithSpan(context.Background(), span)

// 	defer span.Finish()
// 	return nil
// }
