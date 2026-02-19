package storage

import(
	"strconv"
	"internal/logger"
)

type AllowedTypes interface {
	isAllowed()
}

type Gauge float64
func(g Gauge) isAllowed() {}

type Count int64
func(c *Count) isAllowed() {}

func (c *Count) Tick(times Count) {
	*c += times
}

type memStorage[T AllowedTypes] struct {
	l logger.Logger
	ds map[string]T
}

func (s *memStorage[T]) Set(name string, value string) (err error) {
	switch v := any(s.ds).(type) {
		case map[string]Gauge:
			var newValue float64
			if newValue, err = strconv.ParseFloat(value, 64); err == nil {
				v[name] = Gauge(newValue)
				return nil
			}
		case map[string]*Count:
			var newValue int64
			if newValue, err = strconv.ParseInt(value, 10, 64); err == nil {
				if v[name] == nil {
					v[name] = new(Count)
				}
				v[name].Tick(Count(newValue))
				return nil
			}
		}
	return err
}

func (s *memStorage[T]) Get(name string) (any, bool) {
	if value, ok := s.ds[name]; ok {
		return value, ok
	} else {
		return nil, ok
	}
}

func (s *memStorage[T]) Remove(name string) {
	delete(s.ds, name)
}

func NewStorage[T AllowedTypes](l logger.Logger) (*memStorage[T]) {
	return &memStorage[T]{l: l, ds: map[string]T{}}
}
