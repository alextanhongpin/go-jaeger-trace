package main

import (
	"fmt"
	"net/http"

	"github.com/alextanhongpin/go-jaeger-trace/tracer"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
)

const port = ":8080"

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:8080/redirect", 301)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"message": "hello!"}`)
}

func main() {
	t, closer := tracer.New("server", "localhost:5775")
	defer closer.Close()
	opentracing.SetGlobalTracer(t)

	http.HandleFunc("/", index)
	http.HandleFunc("/redirect", redirect)
	fmt.Printf("listening to port *%s. press ctrl + c to cancel", port)
	http.ListenAndServe(port, nethttp.Middleware(opentracing.GlobalTracer(), http.DefaultServeMux))
}
