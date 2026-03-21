package main

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/agent"
	"time"
	"fmt"
	"os/signal"
	"syscall"
	"context"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := run(); err != nil {
		log.Debug().Err(err).Msg("server error")
		panic(err)
	}
}

func run() (err error) {
	err = parseFlags()
	if err != nil { return err }

	ctx, cancelFunc := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelFunc()

	a := agent.NewAgent(flagReportAddr)
	errCh := make(chan error, 2)

	go a.UpdateWorker.Run(ctx, errCh, time.Duration(flagPollInterval) * time.Second)
	go a.ReportWorker.Run(ctx, errCh, time.Duration(flagReportInterval) * time.Second)

	fmt.Println("Agent working, for exit press Ctrl+C")

	for {
		select {
			case <-ctx.Done(): return nil
			case err = <-errCh: return err
		}
	}
}
