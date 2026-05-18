package repository

import(
	"iter"
	"bufio"
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"
	"context"
	"github.com/rs/zerolog/log"
	"maps"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/model"
	"fmt"
)

type memStorage struct {
	storageFile *os.File
	mu sync.Mutex
	ds map[string]models.Metrics
}

func (s *memStorage) Set(ctx context.Context, v models.Metrics) error {
	s.mu.Lock()
	s.ds[v.Key()] = v
	s.mu.Unlock()
	return nil
}

func (s *memStorage) BulkSet(ctx context.Context, vs []models.Metrics) error {
	s.mu.Lock()

	iterVs := func(yield func(string, models.Metrics) bool) {
		for _, v := range vs { if !yield(v.Key(), v) { return } }
	}
	maps.Insert(s.ds, iterVs)
	s.mu.Unlock()
	return nil
}

func (s *memStorage) Get(ctx context.Context, k string) (models.Metrics, error) {
	var err error
	s.mu.Lock()
	value, ok := s.ds[k]
	s.mu.Unlock()
	if !ok {
		err = ErrNotFound
	}
	return value, err
}

func (s *memStorage) Remove(ctx context.Context, k string) error {
	delete(s.ds, k)
	return nil
}

func NewStorage(storageFile *os.File, restore bool) (storage *memStorage, err error) {
	storage = &memStorage{ds: map[string]models.Metrics{}, storageFile: storageFile}
	if restore {
		err = storage.restore()
		if err != nil {
			log.Debug().Err(err).Msg("storage restore error")
			return nil, fmt.Errorf("failed to restore storage from storage file: %w", err)
		}
	}

	return storage, nil
}

func (s *memStorage) All(ctx context.Context) (iter.Seq[models.Metrics], error) {
	return func(yield func(models.Metrics) bool) {
		for _, v := range s.ds {
			if !yield(v) { return }
		}
	}, nil
}

func (s *memStorage) WriteToFile(ctx context.Context, errCh chan error, delay time.Duration) {
	writer := bufio.NewWriter(s.storageFile)

	for {
		select {
			case <-ctx.Done(): return
			case <-time.After(delay):
				s.mu.Lock()
				err := s.storageFile.Truncate(0)
				if err != nil {
					log.Debug().Err(err).Msg("file truncate err")
					errCh <- fmt.Errorf("failed to truncate storage file: %w", err)
					return
				}

				_, err = s.storageFile.Seek(0, io.SeekStart)
				if err != nil {
					log.Debug().Err(err).Msg("file seek err")
					errCh <- fmt.Errorf("failed to seek storage file: %w", err)
					return
				}

				items, err := s.All(ctx)

				if err != nil {
					log.Debug().Err(err).Msg("failed fetch metrics")
					errCh <- fmt.Errorf("failed to marshal metric: %w", err)
					return
				}

				for item := range items {
					data, err := json.Marshal(item)
					if err != nil {
						log.Debug().Err(err).Str("itemKey", item.Key()).Msg("marshal item err")
						errCh <- fmt.Errorf("failed to marshal metric: %w", err)
						return
					}

					_, err = writer.Write(data)
					if err != nil {
						log.Debug().Err(err).Str("data", string(data)).Msg("write item err")
						errCh <- fmt.Errorf("failed to write data to buffer: %w", err)
						return
					}

					err = writer.WriteByte('\n')
					if err != nil {
						log.Debug().Err(err).Msg("write byte err")
						errCh <- fmt.Errorf("failed to write data to buffer: %w", err)
						return
					}

					writer.Flush()
				}
				s.mu.Unlock()
		}
	}
}

func (s *memStorage) restore() (err error) {
	scanner := bufio.NewScanner(s.storageFile)

	for scanner.Scan() {
		var item models.Metrics
		err = json.Unmarshal([]byte(scanner.Text()), &item)
		if err != nil {
			log.Debug().Err(err).Str("item data", scanner.Text()).Msg("unmarshal item err")
			return fmt.Errorf("failed to unmarshal object from storage file: %w", err)
		}
		err := s.Set(context.Background(), item)

		if err != nil {
			return fmt.Errorf("failed to set object to repository: %w", err)
		}
	}

	if err = scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan storage file: %w", err)
	}
	return nil
}

func (s *memStorage) Ping(ctx context.Context) error {
	return nil
}
