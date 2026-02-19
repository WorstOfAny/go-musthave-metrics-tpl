package handler

import(
	"net/http"
	"net/http/httptest"
	"testing"
	"io"
	"github.com/stretchr/testify/assert"
	"github.com/go-chi/chi/v5"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/model"
)


func TestUpdate(t *testing.T) {
	c := NewMetricsController()
	testCases := []struct {
		name string
		method string
		expectedCode int
		response string
		metrics []models.Metrics
		path string
		contentType string
		responseContentType string
	}{
		{
			name: "GET Request",
			path: "/update/gauge/some/1",
			method: http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
			response: "",
			contentType: "text/plain; charset=utf-8",
			responseContentType: "text/plain; charset=utf-8",
		},
		{
			name: "PUT Request",
			path: "/update/gauge/some/1",
			method: http.MethodPut,
			expectedCode: http.StatusMethodNotAllowed,
			response: "",
			contentType: "text/plain; charset=utf-8",
			responseContentType: "text/plain; charset=utf-8",
		},
		{
			name: "DELETE Request",
			path: "/update/gauge/some/1",
			method: http.MethodDelete,
			expectedCode: http.StatusMethodNotAllowed,
			response: "",
			contentType: "text/plain; charset=utf-8",
			responseContentType: "text/plain; charset=utf-8",
		},
		{
			name: "POST Request Update with invalid content type",
			path: "/update/garbge/some/1",
			method: http.MethodPost,
			expectedCode: http.StatusUnsupportedMediaType,
			response: "",
			contentType: "application/json",
			responseContentType: "text/plain; charset=utf-8",
		},
		{
			name: "POST Request Update with invalid type parameter",
			path: "/update/garbge/some/1",
			method: http.MethodPost,
			expectedCode: http.StatusBadRequest,
			response: "",
			contentType: "text/plain; charset=utf-8",
			responseContentType: "text/plain; charset=utf-8",
		},
		{
			name: "POST Request Update gauge with valid parameters",
			path: "/update/gauge/some/1",
			method: http.MethodPost,
			expectedCode: http.StatusOK,
			response: "",
			contentType: "text/plain; charset=utf-8",
			responseContentType: "text/plain; charset=utf-8",
		},
		{
			name: "Get Request all metrics",
			path: "/",
			method: http.MethodGet,
			expectedCode: http.StatusOK,
			response: "<html><body><p>gauge some: 1</p></body></html>",
			contentType: "text/plain; charset=utf-8",
			responseContentType: "text/html;",
		},
		{
			name: "POST Request Update gauge with invalid varName parameter",
			path: "/update/gauge/1",
			method: http.MethodPost,
			expectedCode: http.StatusNotFound,
			response: "404 page not found\n",
			contentType: "text/plain; charset=utf-8",
			responseContentType: "text/plain; charset=utf-8",
		},
		{
			name: "POST Request Update gauge with invalid varValue parameter",
			path: "/update/gauge/some/_",
			method: http.MethodPost,
			expectedCode: http.StatusBadRequest,
			response: "",
			contentType: "text/plain; charset=utf-8",
			responseContentType: "text/plain; charset=utf-8",
		},
		{
			name: "POST Request Update counter with valid parameters",
			path: "/update/counter/some/1",
			method: http.MethodPost,
			expectedCode: http.StatusOK,
			response: "",
			contentType: "text/plain; charset=utf-8",
			responseContentType: "text/plain; charset=utf-8",
		},
		{
			name: "POST Request Update counter with invalid varName parameter",
			path: "/update/counter/1",
			method: http.MethodPost,
			expectedCode: http.StatusNotFound,
			response: "404 page not found\n",
			contentType: "text/plain; charset=utf-8",
			responseContentType: "text/plain; charset=utf-8",
		},
		{
			name: "POST Request Update counter with invalid varValue parameter",
			path: "/update/counter/some/_",
			method: http.MethodPost,
			expectedCode: http.StatusBadRequest,
			response: "",
			contentType: "text/plain; charset=utf-8",
			responseContentType: "text/plain; charset=utf-8",
		},
		{
			name: "Get gauge metric that exist",
			path: "/value/gauge/someGauge",
			method: http.MethodGet,
			expectedCode: http.StatusOK,
			metrics: []models.Metrics{models.Metrics{ID: "someGauge", MType: models.Gauge, Value: new(float64)}},
			response: "0",
			contentType: "text/plain; charset=utf-8",
			responseContentType: "text/plain; charset=utf-8",
		},
		{
			name: "Get gauge metric that not exist",
			path: "/value/gauge/someGauge1",
			method: http.MethodGet,
			expectedCode: http.StatusNotFound,
			metrics: []models.Metrics{},
			response: "",
			contentType: "text/plain; charset=utf-8",
			responseContentType: "text/plain; charset=utf-8",
		},
		{
			name: "Get myType metric",
			path: "/value/MyGauge/someGauge1",
			method: http.MethodGet,
			expectedCode: http.StatusBadRequest,
			metrics: []models.Metrics{},
			response: "",
			contentType: "text/plain; charset=utf-8",
			responseContentType: "text/plain; charset=utf-8",
		},
		{
			name: "Post with missed metric name",
			path: "/update/gauge//1",
			method: http.MethodPost,
			expectedCode: http.StatusNotFound,
			response: "",
			contentType: "text/plain; charset=utf-8",
			responseContentType: "text/plain; charset=utf-8",
		},
		{
			name: "Post with missed metric value",
			path: "/update/gauge/1//",
			method: http.MethodPost,
			expectedCode: http.StatusBadRequest,
			response: "",
			contentType: "text/plain; charset=utf-8",
			responseContentType: "text/plain; charset=utf-8",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			mux := chi.NewRouter()
			c.ApplyTo(mux)
			r := httptest.NewRequest(tc.method, tc.path, nil)
			r.Header.Set("Content-Type", tc.contentType)
			w := httptest.NewRecorder()
			for _, v := range tc.metrics {

				v.Save()
			}

			mux.ServeHTTP(w, r)

			rBody, _ := io.ReadAll(w.Body)
			rContentType := w.Header().Get("Content-Type")

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tc.response, string(rBody), "Тело ответа не совпадает с ожидаемым")
			assert.Equal(t, tc.responseContentType, rContentType, "Content-Type ответа не совпадает с ожидаемым")
		})
	}
}
