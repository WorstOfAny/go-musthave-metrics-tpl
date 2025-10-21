package main

import(
	"internal/handler"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(`:8080`, handler.UpdateHandler())
}
