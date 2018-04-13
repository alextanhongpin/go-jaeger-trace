package main

import (
	"fmt"
	"net/http"

	"github.com/alextanhongpin/go-jaeger-trace/tracer"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
)

const port = ":8081"

func index(w http.ResponseWriter, r *http.Request) {
	endpoint := fmt.Sprintf("http://localhost%s/redirect", port)
	http.Redirect(w, r, endpoint, 301)
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
