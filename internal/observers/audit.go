package observers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/client"
)

type Observer interface {
	Update(Event)
}

type audit struct {
	client *client.Client
	file   *os.File
	mu     sync.Mutex
}

type Event struct {
	Metrics []string
	TS      time.Time
	IPAddr  string
}

type AuditOptionFunc func(*audit) error

func NewAudit(opts ...AuditOptionFunc) (Observer, error) {
	audit := &audit{}

	for _, opt := range opts {
		err := opt(audit)
		if err != nil {
			return nil, fmt.Errorf("failed apply option for constructor: %w", err)
		}
	}

	return audit, nil
}

func WithURL(url string) AuditOptionFunc {
	return func(a *audit) error {
		client := client.NewClient(url, "")
		a.client = client
		return nil
	}
}

func WithFile(ctx context.Context, filename string) AuditOptionFunc {
	return func(a *audit) error {
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}

		go func() {
			<-ctx.Done()
			file.Close()
		}()
		a.file = file
		return nil
	}
}

func (a *audit) Update(e Event) {
	go a.writeToFile(e)
	go a.send(e)
}

func (a *audit) writeToFile(e Event) {
	if a.file != nil {
		a.mu.Lock()
		defer a.mu.Unlock()

		writer := bufio.NewWriter(a.file)
		data, err := json.Marshal(e)
		if err != nil {
			log.Debug().Err(err).Msg("marshal event err")
			return
		}

		_, err = writer.Write(data)
		if err != nil {
			log.Debug().Err(err).Str("data", string(data)).Msg("write event err")
			return
		}

		err = writer.WriteByte('\n')
		if err != nil {
			log.Debug().Err(err).Msg("write byte err")
			return
		}

		writer.Flush()
	}
}

func (a *audit) send(e Event) {
	if a.client != nil {
		data, err := json.Marshal(e)
		if err != nil {
			log.Debug().Err(err).Msg("marshal error")
			return
		}
		a.client.Post("/", data)
	}
}
