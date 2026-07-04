package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type FailReason string

const (
	Counter                 = "counter"
	Gauge                   = "gauge"
	EmptyID      FailReason = "Empty ID"
	EmptyValue   FailReason = "Empty value"
	WrongType    FailReason = "Wrong metric type"
	InvalidFloat FailReason = "Invalid float"
	InvalidInt   FailReason = "Invalid int"
)

type MetricError struct {
	Reason  FailReason
	Message string
	Err     error
}

func (e *MetricError) Error() string {
	return e.Message
}

func (e *MetricError) Cause() FailReason {
	return e.Reason
}

func (e *MetricError) Unwrap() error {
	return e.Err
}

// NOTE: Не усложняем пример, вводя иерархическую вложенность структур.
// Органичиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// что бы отличать значение "0", от не заданного значения
// и соответственно не кодировать в структуру.

type Metrics struct {
	ID    string   `json:"id" db:"id"`
	MType string   `json:"type" db:"mtype"`
	Delta *int64   `json:"delta,omitempty" db:"delta"`
	Value *float64 `json:"value,omitempty" db:"value"`
	Hash  string   `json:"hash,omitempty" db:"hash"`
}

func NewMetric(mtype string, id string, value string) (*Metrics, error) {
	m := &Metrics{ID: id, MType: mtype}
	updErr := m.Update(value)

	_, validErr := m.Valid()
	if validErr != nil {
		return nil, validErr
	}
	if updErr != nil {
		return nil, updErr
	}

	return m, nil
}

func (m *Metrics) Valid() (bool, error) {
	err := &MetricError{}
	var ok bool
	if m.ID == "" {
		err.Reason = EmptyID
		err.Message = "Empty ID forbidden"
		return ok, err
	}
	switch m.MType {
	case Gauge:
		if m.Value == nil {
			err.Reason = EmptyValue
			err.Message = "Empty value forbidden"
			return ok, err
		}
		ok = true
	case Counter:
		if m.Delta == nil {
			err.Reason = EmptyValue
			err.Message = "Empty value forbidden"
			return ok, err
		}
		ok = true
	default:
		err.Reason = WrongType
		err.Message = "Unsupported metric type"
		return ok, err
	}

	return ok, nil
}

func (m *Metrics) Update(value string) error {
	err := &MetricError{}
	switch m.MType {
	case Gauge:
		if newValue, parseErr := strconv.ParseFloat(value, 64); parseErr == nil {
			if m.Value == nil {
				m.Value = &newValue
			} else {
				*m.Value = newValue
			}
		} else {
			err.Reason = InvalidFloat
			err.Message = "Failed parse float from value"
			return err
		}
	case Counter:
		if newValue, parseErr := strconv.ParseInt(value, 10, 64); parseErr == nil {
			if m.Delta == nil {
				m.Delta = &newValue
			} else {
				*m.Delta += newValue
			}
		} else {
			err.Reason = InvalidInt
			err.Message = "Failed parse int from value"
			return err
		}
	default:
		err.Reason = WrongType
		err.Message = "Unsupported type"
		return err
	}

	return nil
}

func (m Metrics) StringValue() (result string) {
	switch m.MType {
	case Gauge:
		if m.Value != nil {
			result = strconv.FormatFloat(*m.Value, 'g', -1, 64)
		}
	case Counter:
		if m.Delta != nil {
			result = strconv.FormatInt(*m.Delta, 10)
		}
	}

	return result
}

func (m Metrics) String() string {
	return fmt.Sprintf("%s %s: %s", m.MType, m.ID, m.StringValue())
}

func (m Metrics) Key() string {
	return m.MType + m.ID
}

func (m *Metrics) UnmarshalJSON(data []byte) error {
	type Alias Metrics

	al := struct{ *Alias }{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &al); err != nil {
		return err
	}

	ok, cause := m.Valid()

	if !ok {
		return cause
	}

	return nil
}
