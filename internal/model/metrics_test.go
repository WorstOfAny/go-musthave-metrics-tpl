package models

import(
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewMetric(t *testing.T) {
	type want struct {
		returnValue *Metrics
	}

	testCases := []struct {
		name string
		metricID string
		metricType string
		want want
	}{
		{
			name: "Test new metric with valid type",
			metricID: "Test",
			metricType: Gauge,
			want: want {
				returnValue: &Metrics{},
			},
		},
		{
			name: "Test new metric with invalid type",
			metricID: "Test2",
			metricType: "myType",
			want: want {
				returnValue: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := NewMetric(tc.metricType, tc.metricID)
			if err == nil {
				assert.IsType(t, tc.want.returnValue, value)
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tc.want.returnValue, value)
				assert.Error(t, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	gaugeMetric1 := Metrics{ID: "gaugeTest1", MType: Gauge}
	gaugeMetric2 := Metrics{ID: "gaugeTest2", MType: Gauge, Value: new(float64)}
	counterMetric1 := Metrics{ID: "counterTest1", MType: Counter}
	counterMetric2 := Metrics{ID: "counterTest2", MType: Counter, Delta: new(int64)}
	badMetric := Metrics{ID: "Test7", MType: "MyType"}

	type want struct {
		metricValue any
		pointerChanged bool
		returnError bool
	}

	testCases := []struct {
		name string
		metric *Metrics
		newValue string
		want want
	}{
		{
			name: "Test update gauge metric with undefined value",
			metric: &gaugeMetric1,
			newValue: "5.0",
			want: want {
				metricValue: float64(5.0),
				pointerChanged: true,
			},
		},
		{
			name: "Test update gauge metric with defined value",
			metric: &gaugeMetric2,
			newValue: "5.0",
			want: want {
				metricValue: float64(5.0),
				pointerChanged: false,
			},
		},
		{
			name: "Test update gauge metric with invalid value",
			metric: &gaugeMetric1,
			newValue: "TEST",
			want: want {
				metricValue: float64(5.0),
				pointerChanged: false,
				returnError: true,
			},
		},
		{
			name: "Test update counter metric with undefined value",
			metric: &counterMetric1,
			newValue: "5",
			want: want {
				metricValue: int64(5),
				pointerChanged: true,
			},
		},
		{
			name: "Test update counter metric with defined value",
			metric: &counterMetric2,
			newValue: "5",
			want: want {
				metricValue: int64(5),
				pointerChanged: false,
			},
		},
		{
			name: "Test update counter metric with invalid value",
			metric: &counterMetric1,
			newValue: "5.8",
			want: want {
				metricValue: int64(5),
				pointerChanged: false,
				returnError: true,
			},
		},
		{
			name: "Test update MyType metric value",
			metric: &badMetric,
			newValue: "5",
			want: want {
				pointerChanged: false,
				returnError: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.metric.MType {
			case Gauge:
				oldPtr := tc.metric.Value
				err := tc.metric.Update(tc.newValue)
				newPtr := tc.metric.Value

				pointerChanged := oldPtr != newPtr

				assert.Equal(t, tc.want.metricValue, *tc.metric.Value)
				assert.Equal(t, tc.want.pointerChanged, pointerChanged)
				if tc.want.returnError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			case Counter:
				oldPtr := tc.metric.Delta
				err := tc.metric.Update(tc.newValue)
				newPtr := tc.metric.Delta

				pointerChanged := oldPtr != newPtr

				assert.Equal(t, tc.want.metricValue, *tc.metric.Delta)
				assert.Equal(t, tc.want.pointerChanged, pointerChanged)
				if tc.want.returnError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			default:
				oldPtrValue := tc.metric.Value
				oldPtrDelta := tc.metric.Delta
				err := tc.metric.Update(tc.newValue)
				newPtrValue := tc.metric.Value
				newPtrDelta := tc.metric.Delta
				pointerChanged := oldPtrValue != newPtrValue || oldPtrDelta != newPtrDelta

				assert.Nil(t, tc.metric.Value)
				assert.Nil(t, tc.metric.Delta)
				assert.Equal(t, tc.want.pointerChanged, pointerChanged)
				assert.Error(t, err)
			}
		})
	}
}

func TestStringValue(t *testing.T) {
	counterMetric := &Metrics{ID: "MyCounter", MType: Counter, Delta: new(int64)}
	*counterMetric.Delta = 100
	gaugeMetric := &Metrics{ID: "MyGauge", MType: Gauge, Value: new(float64)}
	*gaugeMetric.Value = 112.54621
	blankMetric := &Metrics{}
	blankCounterMetric := &Metrics{MType: Counter}
	blankGaugeMetric := &Metrics{MType: Gauge}

	type want struct {
		returnValue string
	}

	testCases := []struct {
		name string
		metric *Metrics
		want want
	}{
		{
			name: "Test counter string value",
			metric: counterMetric,
			want: want {
				returnValue: "100",
			},
		},
		{
			name: "Test gauge string value",
			metric: gaugeMetric,
			want: want {
				returnValue: "112.54621",
			},
		},
		{
			name: "Test blank string value",
			metric: blankMetric,
			want: want {
				returnValue: "",
			},
		},
		{
			name: "Test blank counter string value",
			metric: blankCounterMetric,
			want: want {
				returnValue: "",
			},
		},
		{
			name: "Test blank gauge string value",
			metric: blankGaugeMetric,
			want: want {
				returnValue: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want.returnValue, tc.metric.StringValue())
		})
	}
}

func TestString(t *testing.T) {
	t.Run("Should return string metric representation", func(t *testing.T) {
		m := Metrics{ID: "MyM", MType: Gauge, Value: new(float64)}

		assert.Equal(t, "gauge MyM: 0", m.String())
	})
}
