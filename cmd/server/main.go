package main

import(
	"internal/handler"
	"internal/storage"
	"internal/logger"
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
	repos := map[string]handler.Repository{
		"gauge": storage.NewStorage[storage.Gauge](l),
		"counter": storage.NewStorage[*storage.Count](l),
	}
	c := handler.NewController(l, repos)
	mux := http.NewServeMux()
	c.ApplyTo(mux)
	return http.ListenAndServe(`:8080`, mux)
}
