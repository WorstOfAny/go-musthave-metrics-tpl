package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"

	models "github.com/WorstOfAny/go-musthave-metrics-tpl/internal/model"
)

func TestNewStats(t *testing.T) {
	t.Run("Create Stats object", func(t *testing.T) {
		NewStats()
	})
}

func TestUpdate(t *testing.T) {
	s := NewStats()

	t.Run("Should update object fields", func(t *testing.T) {
		oldPollCount := *s.PollCount
		oldRandomValue := *s.RandomValue
		s.UpdateRT()
		assert.NotEqual(t, oldPollCount, *s.PollCount)
		assert.NotEqual(t, oldRandomValue, *s.RandomValue)
	})
}

func TestAllMetrics(t *testing.T) {
	s := NewStats()

	t.Run("Should iterate through stats metrics", func(t *testing.T) {
		for v := range s.AllMetrics() {
			assert.IsType(t, models.Metrics{}, *v)
		}
		for v := range s.AllMetrics() {
			assert.NotNil(t, v)
			break
		}
	})
}
