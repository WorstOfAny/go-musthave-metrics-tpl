package main

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/logger"
	"net/http"
	"fmt"
)

func LogOutput(message ...any) {
	fmt.Println(message...)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	l := logger.LoggerAdapter(LogOutput)
	c := handler.NewController(l)
	mux := http.NewServeMux()
	c.ApplyTo(mux)
	return http.ListenAndServe(`:8080`, mux)
}
