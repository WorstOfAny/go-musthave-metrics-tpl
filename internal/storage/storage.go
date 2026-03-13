package storage

import "iter"

type memStorage[T any] struct {
	ds map[string]T
}

func (s *memStorage[T]) Set(k string, v T) {
	s.ds[k] = v
}

func (s *memStorage[T]) Get(k string) (T, bool) {
	value, ok := s.ds[k]
	return value, ok
}

func (s *memStorage[T]) Remove(k string) {
	delete(s.ds, k)
}

func NewStorage[T any]() (*memStorage[T]) {
	return &memStorage[T]{ds: map[string]T{}}
}

func (s *memStorage[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s.ds {
			if !yield(v) { return }
		}
	}
}
