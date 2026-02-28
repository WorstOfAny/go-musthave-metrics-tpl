package storage

import(
	"iter"
	"os"
	"bufio"
	"encoding/json"
)

type hasKey interface {
	Key() string
}

type memStorage[T hasKey] struct {
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

func NewStorage[T hasKey]() (*memStorage[T]) {
	return &memStorage[T]{ds: map[string]T{}}
}

func (s *memStorage[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s.ds {
			if !yield(v) { return }
		}
	}
}

func (ms memStorage[T]) WriteToFile(filename string) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	writer := bufio.NewWriter(file)
	if err != nil { return }

	defer file.Close()

	for item := range ms.All() {
		data, err := json.Marshal(item)
		if err != nil { return }
		_, err = writer.Write(data)
		if err != nil { return }
		err = writer.WriteByte('\n')
		if err != nil { return }
		writer.Flush()
	}
}


func (ms *memStorage[T]) RestoreFromFile(filename string) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil { return }
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var item T
		json.Unmarshal([]byte(scanner.Text()), &item)
		ms.Set(item.Key(), item)
	}

	if err = scanner.Err(); err != nil { return }
}
