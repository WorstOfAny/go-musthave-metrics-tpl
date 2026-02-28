package main

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/storage"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/model"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
	"context"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	parseFlags()
	st := storage.NewStorage[*models.Metrics]()

	if flagRestoreStorage {
		st.RestoreFromFile(flagFileStoragePath)
	}

	go func() {
		for {
			select {
				case <-ctx.Done(): return
				case <-time.After(time.Duration(flagStoreInterval) * time.Second): st.WriteToFile(flagFileStoragePath)
			}
		}
	}()
	c := handler.NewMetricsController(st)
	r := chi.NewRouter()
	c.ApplyTo(r)
	srv := &http.Server{Addr: flagRunAddr, Handler: r}
	defer srv.Close()
	return srv.ListenAndServe()
	<-ctx.Done()
	return ctx.Err()
}
