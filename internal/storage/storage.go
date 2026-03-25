package storage

import(
	"iter"
	"os"
	"bufio"
	"encoding/json"
	"time"
	"context"
	"io"
	"sync"
	"github.com/rs/zerolog/log"
)

type hasKey interface {
	Key() string
}

type memStorage[T hasKey] struct {
	storageFile *os.File
	mu sync.Mutex
	ds map[string]T
}

func (s *memStorage[T]) Set(k string, v T) {
	s.mu.Lock()
	s.ds[k] = v
	s.mu.Unlock()
}

func (s *memStorage[T]) Get(k string) (T, bool) {
	s.mu.Lock()
	value, ok := s.ds[k]
	s.mu.Unlock()
	return value, ok
}

func (s *memStorage[T]) Remove(k string) {
	delete(s.ds, k)
}

func NewStorage[T hasKey](storageFile *os.File, restore bool) (storage *memStorage[T], err error) {
	storage = &memStorage[T]{ds: map[string]T{}, storageFile: storageFile}
	if restore {
		err = storage.restore()
		if err != nil {
			log.Debug().Str("storage restore error", err.Error()).Msg("")
			return nil, err
		}
	}

	return storage, nil
}

func (s *memStorage[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s.ds {
			if !yield(v) { return }
		}
	}
}

func (ms *memStorage[T]) WriteToFile(ctx context.Context, errCh chan error, delay time.Duration) {
	writer := bufio.NewWriter(ms.storageFile)

	for {
		select {
			case <-ctx.Done(): return
			case <-time.After(delay):
				ms.mu.Lock()
				err := ms.storageFile.Truncate(0)
				if err != nil {
					log.Debug().Err(err).Msg("file truncate err")
					errCh <- err
					return
				}

				_, err = ms.storageFile.Seek(0, io.SeekStart)
				if err != nil {
					log.Debug().Err(err).Msg("file seek err")
					errCh <- err
					return
				}

				for item := range ms.All() {
					data, err := json.Marshal(item)
					if err != nil {
						log.Debug().Err(err).Str("itemKey", item.Key()).Msg("marshal item err")
						errCh <- err
						return
					}

					_, err = writer.Write(data)
					if err != nil {
						log.Debug().Err(err).Str("data", string(data)).Msg("write item err")
						errCh <- err
						return
					}

					err = writer.WriteByte('\n')
					if err != nil {
						log.Debug().Err(err).Msg("write byte err")
						errCh <- err
						return
					}

					writer.Flush()
				}
				ms.mu.Unlock()
		}
	}
}

func (ms *memStorage[T]) restore() (err error) {
	scanner := bufio.NewScanner(ms.storageFile)

	for scanner.Scan() {
		var item T
		err = json.Unmarshal([]byte(scanner.Text()), &item)
		if err != nil {
			log.Debug().Err(err).Str("item data", scanner.Text()).Msg("unmarshal item err")
			return err
		}
		ms.Set(item.Key(), item)
	}

	if err = scanner.Err(); err != nil { return err }
	return nil
}
