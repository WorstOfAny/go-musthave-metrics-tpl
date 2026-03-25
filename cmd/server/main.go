package main

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/storage"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/model"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
	"context"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"github.com/rs/zerolog/log"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	if err := run(); err != nil {
		log.Debug().Err(err).Msg("server run error")
		panic(err)
	}
}

func run() (err error) {
	err = parseFlags()
	if err != nil { return err }

	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil { return err }
	defer db.Close()

	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil { return err }
	defer file.Close()

	ctx, cancelFunc := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelFunc()

	st, err := storage.NewStorage[*models.Metrics](file, cfg.RestoreStorage)
	if err != nil { return err }

	errCh := make(chan error, 2)
	go st.WriteToFile(ctx, errCh, time.Duration(cfg.StoreInterval) * time.Second)

	c := handler.NewMetricsController(st, db)
	r := chi.NewRouter()
	c.ApplyTo(r)

	go runServer(errCh, cfg.RunAddr, r)

	for {
		select {
			case <-ctx.Done(): return nil
			case err = <-errCh: return err
		}
	}
}

func runServer(errCh chan error, addr string, r *chi.Mux) {
	fmt.Println("Server working, for exit press Ctrl+C")
	srv := &http.Server{Addr: addr, Handler: r}
	defer srv.Close()
	errCh <- srv.ListenAndServe()
}
