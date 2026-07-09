package resetable_pool

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	models "github.com/WorstOfAny/go-musthave-metrics-tpl/internal/model"
)

func TestNew(t *testing.T) {
	p := New[*bytes.Buffer](nil)
	t.Run("should create pool instance", func(t *testing.T) {
		assert.NotNil(t, p)
	})
}
func TestGet(t *testing.T) {
	p := New[*bytes.Buffer](func() *bytes.Buffer {
		return new(bytes.Buffer)
	})
	t.Run("should return *bytes.Buffer", func(t *testing.T) {
		val := p.Get()
		assert.IsType(t, &bytes.Buffer{}, val)
		val2 := p.Get()
		assert.IsType(t, &bytes.Buffer{}, val2)
		assert.NotSame(t, val, val2)
	})
}
func TestPut(t *testing.T) {
	p := New[*models.Metrics](func() *models.Metrics {
		return new(models.Metrics)
	})
	t.Run("should put *models.Metrics into pool", func(t *testing.T) {
		val := p.Get()
		assert.IsType(t, &models.Metrics{}, val)
		paddr := &val
		p.Put(val)
		val = p.Get()
		val.MType = "gauge"
		assert.Equal(t, paddr, &val)
		val2 := p.Get()
		assert.NotSame(t, val, val2)
		assert.NotEqual(t, val, val2)
	})

}
