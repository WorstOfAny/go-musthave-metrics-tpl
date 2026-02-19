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
	parseFlags()
	c := handler.NewMetricsController()
	r := chi.NewRouter()
	c.ApplyTo(r)

	srv := &http.Server{Addr: flagRunAddr, Handler: r}
	defer srv.Close()
	return srv.ListenAndServe()
}
