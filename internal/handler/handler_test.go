package handler

import(
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"io"
	"github.com/stretchr/testify/assert"
	"github.com/go-chi/chi/v5"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/model"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/repository"
	"math"
	"context"
)

type want struct {
	code int
	contentType string
	body string
}

type bodyCase struct {
	name string
	body string
	want want
}

type pathCase struct {
	name string
	path string
	bodyCases []bodyCase
}

type methodCase struct {
	name string
	method string
	pathCases []pathCase
}

type requestCase struct {
	name string
	contentType string
	body string
	methodCases []methodCase
}


func TestListAll(t *testing.T) {
	testCases := []requestCase {
		requestCase{
			name: "Content-Type: text_plain",
			contentType: "text/plain",
			methodCases: []methodCase{
				methodCase{
					name: "Method: GET",
					method: http.MethodGet,
					pathCases: []pathCase{
						pathCase{
							name: "Index path",
							path: "/",
							bodyCases: []bodyCase{
								bodyCase{
									name: "with empty body",
									want: want{ code: http.StatusOK, contentType: "text/html", body: "<html><body><p>gauge MyGauge: 1.1</p></body></html>" },
								},
							},
						},
					},
				},
			},
		},
		requestCase{
			name: "Content-Type: my_content",
			contentType: "my/content",
			methodCases: []methodCase{
				methodCase{
					name: "Method: GET",
					method: http.MethodGet,
					pathCases: []pathCase{
						pathCase{
							name: "Index path",
							path: "/",
							bodyCases: []bodyCase{
								bodyCase{
									name: "Any",
									want: want{ code: http.StatusUnsupportedMediaType, contentType: "text/plain", body: "" },
								},
							},
						},
					},
				},
			},
		},
	}

	val := 1.1

	mcs := []models.Metrics{
		models.Metrics{ID:"MyGauge", MType: models.Gauge, Value: &val},
	}

	runCases(t, testCases, mcs)
}

func TestShowTextPlain(t *testing.T) {
	testCases := []requestCase {
		requestCase{
			name: "Content-Type: text_plain",
			contentType: "text/plain",
			methodCases: []methodCase{
				methodCase{
					name: "Method: GET",
					method: http.MethodGet,
					pathCases: []pathCase{
						pathCase{
							name: "Show path with valid path parameters",
							path: "/value/gauge/MyGauge",
							bodyCases: []bodyCase{
								bodyCase{
									name: "with empty body",
									want: want{ code: http.StatusOK, contentType: "text/plain", body: "1.1" },
								},
							},
						},
						pathCase{
							name: "Show path with path parameter metricName that doesn't exist",
							path: "/value/gauge/UnknownGauge",
							bodyCases: []bodyCase{
								bodyCase{
									name: "with empty body",
									want: want{ code: http.StatusNotFound, contentType: "text/plain", body: "" },
								},
							},
						},
						pathCase{
							name: "Show path with path parameter metricType that doesn't exist",
							path: "/value/MyType/Unknown",
							bodyCases: []bodyCase{
								bodyCase{
									name: "with empty body",
									want: want{ code: http.StatusNotFound, contentType: "text/plain", body: "" },
								},
							},
						},
					},
				},
			},
		},
	}

	val := 1.1

	mcs := []models.Metrics{
		models.Metrics{ID:"MyGauge", MType: models.Gauge, Value: &val},
	}

	runCases(t, testCases, mcs)
}

func TestShowJSON(t *testing.T) {
	testCases := []requestCase {
		requestCase{
			name: "Content-Type: application_json",
			contentType: "application/json",
			methodCases: []methodCase{
				methodCase{
					name: "Method: POST",
					method: http.MethodPost,
					pathCases: []pathCase{
						pathCase{
							name: "Show path",
							path: "/value",
							bodyCases: []bodyCase{
								bodyCase{
									name: "with valid json body get Gauge",
									body: "{\"id\": \"MyGauge\", \"type\": \"gauge\"}",
									want: want{ code: http.StatusOK, contentType: "application/json", body: "{\"id\":\"MyGauge\",\"type\":\"gauge\",\"value\":1.1}" },
								},
								bodyCase{
									name: "with valid json body get Counter",
									body: "{\"id\": \"MyCounter\", \"type\": \"counter\"}",
									want: want{ code: http.StatusOK, contentType: "application/json", body: "{\"id\":\"MyCounter\",\"type\":\"counter\",\"delta\":5}" },
								},
								bodyCase{
									name: "with valid json body but not existing metric",
									body: "{\"id\": \"SomeGauge\", \"type\": \"gauge\"}",
									want: want{ code: http.StatusNotFound, contentType: "application/json", body: "" },
								},
								bodyCase{
									name: "with valid json body and invalid type",
									body: "{\"id\": \"MyGauge\", \"type\": \"MyType\"}",
									want: want{ code: http.StatusNotFound, contentType: "application/json", body: "" },
								},
								bodyCase{
									name: "with invalid json body",
									body: "{id: \"MyGauge\", type: \"gauge\"}",
									want: want{ code: http.StatusBadRequest, contentType: "application/json", body: "" },
								},
								bodyCase{
									name: "with invalid json body",
									body: "{id: \"MyGauge\", type: \"gauge\"",
									want: want{ code: http.StatusBadRequest, contentType: "application/json", body: "" },
								},
								bodyCase{
									name: "with valid json body, but object can't be serialized",
									body: "{\"id\": \"BrokenGauge\", \"type\": \"gauge\"}",
									want: want{ code: http.StatusInternalServerError, contentType: "application/json", body: "" },
								},
							},
						},
					},
				},
			},
		},
	}

	gaugeVal := 1.1
	counterVal := int64(5)
	brokenVal := math.NaN()

	mcs := []models.Metrics{
		models.Metrics{ID:"MyGauge", MType: models.Gauge, Value: &gaugeVal},
		models.Metrics{ID:"BrokenGauge", MType: models.Gauge, Value: &brokenVal},
		models.Metrics{ID:"MyCounter", MType: models.Counter, Delta: &counterVal},
	}


	runCases(t, testCases, mcs)
}

func TestUpdate(t *testing.T) {
	testCases := []requestCase {
		requestCase{
			name: "Content-Type: text_plain",
			contentType: "text/plain",
			methodCases: []methodCase{
				methodCase{
					name: "Method: POST",
					method: http.MethodPost,
					pathCases: []pathCase{
						pathCase{
							name: "Update path with valid path parameters with old ID",
							path: "/update/gauge/MyGauge/1.1",
							bodyCases: []bodyCase{
								bodyCase{
									name: "with empty body",
									want: want{ code: http.StatusOK, contentType: "text/plain", body: "{\"id\":\"MyGauge\",\"type\":\"gauge\",\"value\":1.1}"  },
								},
							},
						},
						pathCase{
							name: "Update path with valid path parameters with new ID",
							path: "/update/gauge/newGauge/1.1",
							bodyCases: []bodyCase{
								bodyCase{
									name: "with empty body",
									want: want{ code: http.StatusOK, contentType: "text/plain", body: "{\"id\":\"newGauge\",\"type\":\"gauge\",\"value\":1.1}" },
								},
							},
						},
						pathCase{
							name: "Update path with invalid path parameter metricValue",
							path: "/update/gauge/newGauge/one+dot+one",
							bodyCases: []bodyCase{
								bodyCase{
									name: "with empty body",
									want: want{ code: http.StatusBadRequest, contentType: "text/plain", body: "" },
								},
							},
						},
						pathCase{
							name: "Update path with missed path parameter metricValue",
							path: "/update/gauge/newGauge//",
							bodyCases: []bodyCase{
								bodyCase{
									name: "with empty body",
									want: want{ code: http.StatusBadRequest, contentType: "text/plain", body: "" },
								},
							},
						},
						pathCase{
							name: "Update path with invalid path parameter metricName",
							path: "/update/gauge//one+dot+one/",
							bodyCases: []bodyCase{
								bodyCase{
									name: "with empty body",
									want: want{ code: http.StatusNotFound, contentType: "text/plain", body: "" },
								},
							},
						},
						pathCase{
							name: "Update path with invalid path parameter metricType",
							path: "/update/MyType/myObject/one+dot+one/",
							bodyCases: []bodyCase{
								bodyCase{
									name: "with empty body",
									want: want{ code: http.StatusBadRequest, contentType: "text/plain", body: "" },
								},
							},
						},
					},
				},
			},
		},
		requestCase{
			name: "Request application/json",
			contentType: "application/json",
			methodCases: []methodCase{
				methodCase{
					name: "method POST",
					method: http.MethodPost,
					pathCases: []pathCase{
						pathCase{
							name: "Update path",
							path: "/update",
							bodyCases: []bodyCase{
								bodyCase{
									name: "with valid json body update Gauge",
									body: "{\"id\": \"MyGauge\", \"type\": \"gauge\", \"value\": 1.4}",
									want: want{ code: http.StatusOK, contentType: "application/json", body: "{\"id\":\"MyGauge\",\"type\":\"gauge\",\"value\":1.4}" },
								},
								bodyCase{
									name: "with valid json body update Counter",
									body: "{\"id\": \"MyCounter\", \"type\": \"counter\", \"delta\": 5}",
									want: want{ code: http.StatusOK, contentType: "application/json", body: "{\"id\":\"MyCounter\",\"type\":\"counter\",\"delta\":5}" },
								},
								bodyCase{
									name: "with invalid json body",
									body: "{id: \"MyGauge\", type: \"gauge\", \"value\": 1.4",
									want: want{ code: http.StatusBadRequest, contentType: "application/json", body: "" },
								},
							},
						},
					},
				},
			},
		},
		requestCase{
			name: "Content-Type: my_content",
			contentType: "my/content",
			methodCases: []methodCase{
				methodCase{
					name: "Method: POST",
					method: http.MethodPost,
					pathCases: []pathCase{
						pathCase{
							name: "Update path",
							path: "/update",
							bodyCases: []bodyCase{
								bodyCase{
									name: "with valid json body",
									body: "{\"id\": \"MyGauge\", \"type\": \"gauge\", \"value\": 1.4}",
									want: want{ code: http.StatusUnsupportedMediaType, contentType: "", body: "" },
								},
							},
						},
					},
				},
			},
		},
	}


	mcs := []models.Metrics{
		models.Metrics{ID:"MyGauge", MType: models.Gauge},
		models.Metrics{ID:"MyCounter", MType: models.Counter},
	}

	runCases(t, testCases, mcs)

}

func runCases(t *testing.T, cases []requestCase, metrics []models.Metrics) {
	for _, tc := range cases {
		testConfig := repository.Config{StoreInterval: 10, FileStoragePath: "test_store.json", RestoreStorage: false, DatabaseDSN: ""}
		errCh := make(chan error, 2)
		ctx := context.Background()
		repo, err := repository.NewRepository(ctx, errCh, testConfig)

		if err != nil { t.Fatal(err) }

		c := NewMetricsController(repo)
		for _, v := range metrics { c.storage.Set(ctx, v) }
		mux := chi.NewRouter()
		c.ApplyTo(mux)
		t.Run(tc.name, func(t *testing.T) {
			for _, mtc := range tc.methodCases {
				t.Run(mtc.name, func(t *testing.T){
					for _, ptc := range mtc.pathCases {
						t.Run(ptc.name, func(t *testing.T) {
							for _, btc := range ptc.bodyCases {
								t.Run(btc.name, func(t *testing.T) {
									r := httptest.NewRequest(mtc.method, ptc.path, strings.NewReader(btc.body))
									r.Header.Set("Content-Type", tc.contentType)
									w := httptest.NewRecorder()
									mux.ServeHTTP(w, r)
									rBody, _ := io.ReadAll(w.Body)
									rContentType := w.Header().Get("Content-Type")
									assert.Equal(t, btc.want.code, w.Code, "Код ответа не совпадает с ожидаемым")
									assert.Equal(t, btc.want.body, string(rBody), "Тело ответа не совпадает с ожидаемым")
									assert.Equal(t, true, strings.HasPrefix(rContentType, btc.want.contentType), "Content-Type ответа не совпадает с ожидаемым")
								})
							}
						})
					}
				})
			}
		})
	}
}
