package models

import(
	"strconv"
	"fmt"
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

func NewMetric(mtype string, id string) (*Metrics, error) {
	switch mtype {
		case Gauge, Counter:
			return &Metrics{ID: id, MType: mtype }, nil
		default:
			return nil, fmt.Errorf("unknown type")
	}
}

func (m *Metrics) Update(value string) error {
	var err error
	switch m.MType {
		case Gauge:
			if newValue, parseErr := strconv.ParseFloat(value, 64); parseErr == nil {
				if m.Value == nil {
					m.Value = &newValue
				} else {
					*m.Value = newValue
				}
			} else {
				err = parseErr
			}
		case Counter:
			if newValue, parseErr := strconv.ParseInt(value, 10, 64); parseErr == nil {
				if m.Delta == nil {
					m.Delta = &newValue
				} else {
					*m.Delta += newValue
				}
			} else {
				err = parseErr
			}
		default:
			err = fmt.Errorf("unknown type")
	}

	return err
}

func (m *Metrics) StringValue() (result string) {
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

func (m *Metrics) String() string {
	return fmt.Sprintf("%s %s: %s", m.MType, m.ID, m.StringValue())
}

func (m *Metrics) Key() string {
	return m.MType + m.ID
}
