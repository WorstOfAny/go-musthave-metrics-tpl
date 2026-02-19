package main

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	c := handler.NewMetricsController()
	r := chi.NewRouter()
	c.ApplyTo(r)
	return http.ListenAndServe(`:8080`, r)
}
