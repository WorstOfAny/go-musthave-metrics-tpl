package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/observers"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/repository"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Build version: %s\n", buildVersion)
	fmt.Fprintf(&buf, "Build date: %s\n", buildDate)
	fmt.Fprintf(&buf, "Build commit: %s\n", buildCommit)
	os.Stdout.Write(buf.Bytes())

	cfg := &config{}
	err := parseFlags(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to parse flags: %w", err))
	}
	if err := run(cfg); err != nil {

		if errors.Is(err, context.Canceled) {
			log.Info().Err(err).Msg("server gracefully shutted down")
			return
		}
		if errors.Is(err, http.ErrServerClosed) {
			log.Info().Err(err).Msg("server gracefully shutted down")
			return
		}
		log.Debug().Err(err).Msg("server run error")
		panic(err)
	}
}

func run(cfg *config) (err error) {
	ctx, cancelFunc := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelFunc()

	repo, err := repository.NewRepository(ctx, cfg.RepoConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}
	privateBytes, err := os.ReadFile(cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("failed read private key")
	}

	c := handler.NewMetricsController(ctx, repo, cfg.Key, privateBytes)
	var auditOpts []observers.AuditOptionFunc
	if cfg.AuditFile != "" {
		auditOpts = append(auditOpts, observers.WithFile(ctx, cfg.AuditFile))
	}
	if cfg.AuditURL != "" {
		auditOpts = append(auditOpts, observers.WithURL(cfg.AuditURL))
	}
	audit, err := observers.NewAudit(auditOpts...)

	if err != nil {
		return fmt.Errorf("failed to initialize audit observer: %w", err)
	}
	c.Register(audit)
	r := chi.NewRouter()
	c.ApplyTo(r)

	err = runServer(ctx, cfg.RunAddr, r)
	return err
}

func runServer(ctx context.Context, addr string, r *chi.Mux) error {
	log.Info().Msgf("server run on: %v", addr)
	srv := &http.Server{Addr: addr, Handler: r}
	defer srv.Close()

	go func() {
		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Info().Err(err).Msg("server forced to shutdown")
		}
	}()
	return srv.ListenAndServe()
}
