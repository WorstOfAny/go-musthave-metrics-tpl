package storage

import(
	"internal/logger"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewStorage (t *testing.T) {
	t.Run("Create Gauge Storage", func(t *testing.T) {
		s := NewStorage[Gauge](nil)

		if s.ds == nil {
			t.Fatal("Map didn't initialized")
		}

		s.ds["myGauge"] = *new(Gauge)
	})

	t.Run("Create Counter Storage", func(t *testing.T) {
		s := NewStorage[*Count](nil)

		if s.ds == nil {
			t.Fatal("Map didn't initialized")
		}

		s.ds["myCounter"] = new(Count)
	})
}

func TestGet(t *testing.T) {
	type want struct {
		returnValue any
		valueExist bool
	}

	storage := map[string]interface{Set(string, string) error; Get(string) (any, bool)}{
		"gauge": NewStorage[Gauge](logger.LoggerAdapter(t.Log)),
		"counter": NewStorage[*Count](logger.LoggerAdapter(t.Log)),
	}

	storage["gauge"].Set("myGauge", "1.1")
	storage["counter"].Set("myCounter", "2")
	tests := []struct {
		name string
		metricType string
		metricName string
		want want
	}{
		{
			name: "Test Gauge that exist",
			metricType: "gauge",
			metricName: "myGauge",
			want: want{returnValue: Gauge(1.1), valueExist: true},
		},
		{
			name: "Test Gauge that not exist",
			metricType: "gauge",
			metricName: "testGauge",
			want: want{returnValue: any(nil), valueExist: false},
		},
		{
			name: "Test Counter that exist",
			metricType: "counter",
			metricName: "myCounter",
			want: want{returnValue: Count(2), valueExist: true},
		},
		{
			name: "Test Counter that not exist",
			metricType: "counter",
			metricName: "testCounter",
			want: want{returnValue: any(nil), valueExist: false},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			value, ok := storage[test.metricType].Get(test.metricName)


			assert.Equal(t, test.want.valueExist, ok)

			switch value.(type) {
				case *Count:
					assert.Equal(t, test.want.returnValue, *value.(*Count))
				default:
					assert.Equal(t, test.want.returnValue, value)
			}
		})
	}
}

func TestSet(t *testing.T) {
	type want struct {
		storedValue any
		returnError bool
	}

	storage := map[string]interface{Set(string, string) error; Get(string) (any, bool)}{
		"gauge": NewStorage[Gauge](logger.LoggerAdapter(t.Log)),
		"counter": NewStorage[*Count](logger.LoggerAdapter(t.Log)),
	}
	tests := []struct {
		name string
		metricType string
		metricName string
		metricValue string
		want want
	}{
		{
			name: "Test Gauge",
			metricType: "gauge",
			metricName: "testGauge",
			metricValue: "100.0",
			want: want{storedValue: Gauge(100.0)},
		},
		{
			name: "Test Gauge with same metric name",
			metricType: "gauge",
			metricName: "testGauge",
			metricValue: "110.0",
			want: want{storedValue: Gauge(110.0)},
		},
		{
			name: "Test Gauge with wrong type",
			metricType: "gauge",
			metricName: "testGauge",
			metricValue: "sss",
			want: want{storedValue: Gauge(110.0), returnError: true},
		},
		{
			name: "Test Counter",
			metricType: "counter",
			metricName: "testCounter",
			metricValue: "5",
			want: want{storedValue: Count(5)},
		},
		{
			name: "Test Counter with same metric name",
			metricType: "counter",
			metricName: "testCounter",
			metricValue: "6",
			want: want{storedValue: Count(11)},
		},
		{
			name: "Test Counter with float type",
			metricType: "counter",
			metricName: "testCounter",
			metricValue: "5.0",
			want: want{storedValue: Count(11), returnError: true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			err := storage[test.metricType].Set(test.metricName, test.metricValue)
			if test.want.returnError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, err, nil)
			}

			value, _ := storage[test.metricType].Get(test.metricName)

			switch value.(type) {
				case *Count:
					assert.Equal(t, test.want.storedValue, *value.(*Count))
				default:
					assert.Equal(t, test.want.storedValue, value)
			}
		})
	}


	// t.Fatal("not implemented")
}

func TestRemove(t *testing.T) {
	storage := map[string]interface{Set(string, string) error; Get(string) (any, bool); Remove(string)}{
		"gauge": NewStorage[Gauge](logger.LoggerAdapter(t.Log)),
		"counter": NewStorage[*Count](logger.LoggerAdapter(t.Log)),
	}

	storage["gauge"].Set("testGauge", "2.1")
	storage["counter"].Set("testCounter", "2")
	tests := []struct {
		name string
		metricType string
		metricName string
	}{
		{
			name: "Test Gauge Remove",
			metricType: "gauge",
			metricName: "testGauge",
		},
		{
			name: "Test Counter remove",
			metricType: "counter",
			metricName: "testCounter",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			storage[test.metricType].Remove(test.metricName)
			_, ok := storage[test.metricType].Get(test.metricName)
			assert.Equal(t, false, ok)
		})
	}
}

func TestTick(t *testing.T) {
	t.Run("Test Tick", func(t *testing.T) {
		c := new(Count)
		c.Tick(Count(2))
		assert.Equal(t, Count(2), *c)
	})
}

func TestIsAllowed(t *testing.T) {
	t.Run("Gauge obj is allowed obj", func(t *testing.T) {
		g := any(new(Gauge))

		if _, ok := g.(interface{isAllowed()}); !ok {
			t.Error("Gauge не реализует isAllowed()")
		}
	})
	t.Run("Count pointer is allowed obj", func(t *testing.T) {
		c := any(new(Count))

		if _, ok := c.(interface{isAllowed()}); !ok {
			t.Error("Count не реализует isAllowed()")
		}
	})
}
