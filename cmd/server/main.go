package main

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/repository"
	"github.com/go-chi/chi/v5"
	"net/http"
	"context"
	"os/signal"
	"syscall"
	"fmt"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := &config{}
	err := parseFlags(cfg)
	if err != nil { panic(fmt.Errorf("failed to parse flags: %w", err)) }
	if err := run(cfg); err != nil {
		log.Debug().Err(err).Msg("server run error")
		panic(err)
	}
}

func run(cfg *config) (err error) {
	ctx, cancelFunc := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelFunc()

	errCh := make(chan error, 2)

	repo, err := repository.NewRepository(ctx, errCh, cfg.RepoConfig)
	if err != nil { return fmt.Errorf("failed to initialize repository: %w", err) }

	c := handler.NewMetricsController(repo, cfg.Key)
	r := chi.NewRouter()
	c.ApplyTo(r)

	go runServer(errCh, cfg.RunAddr, r)

	for {
		select {
			case <-ctx.Done(): return nil
			case err = <-errCh: return fmt.Errorf("server error: %w", err)
		}
	}
}

func runServer(errCh chan error, addr string, r *chi.Mux) {
	fmt.Println("Server working, for exit press Ctrl+C")
	srv := &http.Server{Addr: addr, Handler: r}
	defer srv.Close()
	errCh <- srv.ListenAndServe()
}
