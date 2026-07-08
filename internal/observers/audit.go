package observers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/client"
)

// Observer интерфейс наблюдателя
type Observer interface {
	Update([]byte)
}

// generate:reset

type audit struct {
	client *client.Client
	file   *os.File
	mu     sync.Mutex
}

// NewAudit конструктор для создания наблюдателя типа audit
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

// AuditOptionFunc тип для функциональных опций конструктора
type AuditOptionFunc func(*audit) error

// WithURL функциональная опция для установки URL, куда будут посылаться события
func WithURL(url string) AuditOptionFunc {
	return func(a *audit) error {
		client := client.NewClient(url, "", []byte{})
		a.client = client
		return nil
	}
}

// WithFile функциональная опция для установки пути к файлу, куда будут записываться события
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

// Update отправка данных при возникновении события
func (a *audit) Update(e []byte) {
	go a.writeToFile(e)
	go a.send(e)
}

func (a *audit) writeToFile(data []byte) {
	if a.file != nil {
		a.mu.Lock()
		defer a.mu.Unlock()

		writer := bufio.NewWriter(a.file)
		_, err := writer.Write(data)
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

func (a *audit) send(data []byte) {
	if a.client != nil {
		a.client.Post("/", data)
	}
}
