package main

import (
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"

	"github.com/uber/jaeger-client-go"
	jaegerClientConfig "github.com/uber/jaeger-client-go/config"
)

var tracer opentracing.Tracer

// var closer io.Closer

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
		"your_service_name",
		jaegerClientConfig.Logger(jaeger.StdLogger),
	)

	defer closer.Close()
	testTracer()        // Not working
	testTracer1(tracer) // working
	testTracer2()       // Not working

	http.ListenAndServe(":8080", nil)
}

// Not working
func testTracer() {
	tracer, closer := jaeger.NewTracer("DOOP",
		jaeger.NewConstSampler(true),
		jaeger.NewNullReporter())
	defer closer.Close()

	span := tracer.StartSpan("hello")
	span.SetOperationName("s2")
	span.LogEvent("hello")
	span.SetBaggageItem("Some_Key", "12345")
	span.SetBaggageItem("Some-other-key", "42")
	span.SetTag(zipkincore.HTTP_PATH, struct{ name string }{"ad"})
	// span.SetTag(zipkincore.HTTP_HOST, struct{name string, car string}{"asd", "asd"})
	defer span.Finish()
	span.LogEvent("hello")
}

// Working
func testTracer1(tracer interface{}) interface{} {
	t := tracer.(*jaeger.Tracer)
	span := t.StartSpan("new_span")
	defer span.Finish()
	span.SetOperationName("s2")
	// span.LogFields(log.String("ds", "asd"))
	span.LogEvent("hello")
	span.SetTag(zipkincore.HTTP_PATH, struct{ name string }{"ad"})

	span.SetBaggageItem("Some_Key", "12345")
	span.SetBaggageItem("Some-other-key", "42")
	return nil
}

// Not working
func testTracer2() interface{} {
	span := opentracing.StartSpan("operation_name")

	defer span.Finish()
	span.SetOperationName("s2")
	span.LogFields(log.String("ds", "asd"))
	span.SetTag(zipkincore.HTTP_PATH, struct{ name string }{"ad"})
	span.LogEvent("hello")

	span.SetBaggageItem("Some_Key", "12345")
	span.SetBaggageItem("Some-other-key", "42")
	return nil
}
