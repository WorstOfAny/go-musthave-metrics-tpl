package main

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	c := handler.NewController()
	mux := http.NewServeMux()
	c.ApplyTo(mux)
	return http.ListenAndServe(`:8080`, mux)
}
