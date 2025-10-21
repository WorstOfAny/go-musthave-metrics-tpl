package handler

import(
	"net/http"
	"net/http/httptest"
	"internal/storage"
	"internal/logger"
	"testing"
	"io"
	"github.com/stretchr/testify/assert"
)


func TestUpdate(t *testing.T) {
	l := logger.LoggerAdapter(t.Log)
	repositories := map[string]Repository{
		"gauge": storage.NewStorage[storage.Gauge](l),
		"counter": storage.NewStorage[*storage.Count](l),
	}

	c := NewController(logger.LoggerAdapter(t.Log), repositories)
	testCases := []struct {
		name string
		method string
		expectedCode int
		response string
		path string
		contentType string
	}{
		{ name: "GET Request", path: "/update/gauge/some/1", method: http.MethodGet, expectedCode: http.StatusMethodNotAllowed, response: ""},
		{ name: "PUT Request", path: "/update/gauge/some/1", method: http.MethodPut, expectedCode: http.StatusMethodNotAllowed, response: ""},
		{ name: "DELETE Request", path: "/update/gauge/some/1", method: http.MethodDelete, expectedCode: http.StatusMethodNotAllowed, response: ""},
		{ name: "POST Request Update with invalid content type", path: "/update/garbge/some/1", method: http.MethodPost, expectedCode: http.StatusUnsupportedMediaType, response: "", contentType: "application/json"},
		{ name: "POST Request Update with invalid type parameter", path: "/update/garbge/some/1", method: http.MethodPost, expectedCode: http.StatusBadRequest, response: "", contentType: "text/plain"},
		{ name: "POST Request Update gauge with valid parameters", path: "/update/gauge/some/1", method: http.MethodPost, expectedCode: http.StatusOK, response: "", contentType: "text/plain"},
		{ name: "POST Request Update gauge with invalid varName parameter", path: "/update/gauge/1", method: http.MethodPost, expectedCode: http.StatusNotFound, response: "404 page not found\n", contentType: "text/plain"},
		{ name: "POST Request Update gauge with invalid varValue parameter", path: "/update/gauge/some/_", method: http.MethodPost, expectedCode: http.StatusBadRequest, response: "", contentType: "text/plain"},
		{ name: "POST Request Update counter with valid parameters", path: "/update/counter/some/1", method: http.MethodPost, expectedCode: http.StatusOK, response: "", contentType: "text/plain"},
		{ name: "POST Request Update counter with invalid varName parameter", path: "/update/counter/1", method: http.MethodPost, expectedCode: http.StatusNotFound, response: "404 page not found\n", contentType: "text/plain"},
		{ name: "POST Request Update counter with invalid varValue parameter", path: "/update/counter/some/_", method: http.MethodPost, expectedCode: http.StatusBadRequest, response: "", contentType: "text/plain"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			mux := http.NewServeMux()
			c.ApplyTo(mux)
			r := httptest.NewRequest(tc.method, tc.path, nil)
			r.Header.Set("Content-Type", tc.contentType)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, r)

			rBody, _ := io.ReadAll(w.Body)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tc.response, string(rBody), "Тело ответа не совпадает с ожидаемым")
		})
	}
}
