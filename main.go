package main

import (
	"net/http"
	"time"

	"github.com/alextanhongpin/go-jaeger-trace/tracer"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
)

const jaegerServiceName = "some_service"
const jaegerHostPort = "localhost:5775"

func main() {

	tracer, closer := tracer.New(jaegerServiceName, jaegerHostPort)
	defer closer.Close()
	tracer2(testTracer1(tracer)) // working
	// testTracer()

	// Required so that the tags are sent safely
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
