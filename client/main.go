package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"

	"github.com/opentracing-contrib/go-stdlib/nethttp"

	"github.com/alextanhongpin/go-jaeger-trace/tracer"
	"github.com/opentracing/opentracing-go"
	tlog "github.com/opentracing/opentracing-go/log"
)

var t opentracing.Tracer
var closer io.Closer

func main() {
	t, closer = tracer.New("client", "localhost:5775")
	defer closer.Close()
	opentracing.SetGlobalTracer(t)

	// ctx := context.Background()
	// askGoogle(ctx)
	runClient(t)
}

type clientTrace struct {
	span opentracing.Span
}

func (t *clientTrace) dnsStart(info httptrace.DNSStartInfo) {
	// t.span.LogKV(
	// 	tlog.String("event", "DNS start"),
	// 	tlog.Object("host", info.Host),
	// )
	t.span.LogEvent("DNS start")
	t.span.LogFields(tlog.String("host", info.Host))

}

func (t *clientTrace) dnsDone(httptrace.DNSDoneInfo) {
	t.span.LogFields(tlog.String("event", "DNS done"))
}

func NewClientTrace(span opentracing.Span) *httptrace.ClientTrace {
	trace := &clientTrace{span: span}
	return &httptrace.ClientTrace{
		DNSStart: trace.dnsStart,
		DNSDone:  trace.dnsDone,
	}
}

func runClient(tracer opentracing.Tracer) {
	// nethttp.Transport from go-stdlib will do the tracing
	c := &http.Client{Transport: &nethttp.Transport{}}

	// Create a top-level span to represent full work of the client
	span := tracer.StartSpan("client")
	span.SetTag("hello", "client")
	defer span.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}

	req = req.WithContext(ctx)
	req, ht := nethttp.TraceRequest(tracer, req)
	defer ht.Finish()

	res, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(body))
}

func askGoogle(ctx context.Context) {
	var parentCtx opentracing.SpanContext
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan != nil {
		parentCtx = parentSpan.Context()
	}

	// Start a new span to wrap HTTP request
	span := t.StartSpan("ask google", opentracing.ChildOf(parentCtx))
	defer span.Finish()

	// Make the Span current in the context
	ctx = opentracing.ContextWithSpan(ctx, span)
	req, err := http.NewRequest("GET", "http://google.com", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Attach ClientTrace to the Context, and Context to request
	trace := NewClientTrace(span)
	ctx = httptrace.WithClientTrace(ctx, trace)
	req = req.WithContext(ctx)

	// Execute the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
}
