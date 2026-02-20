package models

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/storage"
	"strconv"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

// NOTE: Не усложняем пример, вводя иерархическую вложенность структур.
// Органичиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// что бы отличать значение "0", от не заданного значения
// и соответственно не кодировать в структуру.

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

var metricsStorage = storage.NewStorage[*Metrics](nil)

func NewMetric(mtype string, id string) (*Metrics, bool) {
	switch mtype {
		case Gauge, Counter:
			return &Metrics{ID: id, MType: mtype }, true
		default:
			return nil, false
	}
}

func FindMetric(mtype string, id string) (m *Metrics, ok bool) {
	m, ok = metricsStorage.Get(mtype + id)
	return m, ok
}

func (m *Metrics) Save() {
	metricsStorage.Set(m.MType + m.ID, m)
}

func (m *Metrics) Update(value string) (success bool) {
	switch m.MType {
		case Gauge:
			if newValue, err := strconv.ParseFloat(value, 64); err == nil {
				if m.Value == nil {
					m.Value = &newValue
				} else {
					*m.Value = newValue
				}
				success = true
			}
		case Counter:
			if newValue, err := strconv.ParseInt(value, 10, 64); err == nil {
				if m.Delta == nil {
					m.Delta = &newValue
				} else {
					*m.Delta += newValue
				}
				success = true
			}
	}

	return success
}

func (m *Metrics) StringValue() string {
	var result string

	switch m.MType {
	case Gauge:
		if m.Value == nil { return result }
		result = strconv.FormatFloat(*m.Value, 'g', -1, 64)
	case Counter:
		if m.Delta == nil { return result }
		result = strconv.FormatInt(*m.Delta, 10)
	}

	return result
}
