package handler

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/model"
	"github.com/go-chi/chi/v5"
	//"github.com/go-chi/chi/v5/middleware"
	//"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"context"
	"fmt"
	"iter"
	"time"
)

type metricsStorage interface {
 Set(string, *models.Metrics)
 Get(string) (*models.Metrics, bool)
 Remove(string)
 All() iter.Seq[*models.Metrics]
}

type metricsController struct {
	storage metricsStorage
}

func NewMetricsController(storage metricsStorage) *metricsController {
	return &metricsController{storage: storage}
}

type responseData struct {
	status int
	size int
}
type responseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (wr *responseWriter) Write(b []byte) (int, error) {
	size, err := wr.ResponseWriter.Write(b)
	wr.responseData.size += size
	return size, err
}

func (wr *responseWriter) WriteHeader(statusCode int) {
	wr.ResponseWriter.WriteHeader(statusCode)
	wr.responseData.status = statusCode
}

type ctxKey string
const(
	metricKey ctxKey = "metric"
	metricTypeKey ctxKey = "metricType"
	metricNameKey ctxKey = "metricName"
	metricValueKey ctxKey = "metricValue"
)

func (c *metricsController) ApplyTo(mux chi.Router) {
	mux.Use(logRequest)
	mux.Use(textPlainTypeSet)
	mux.Get("/", c.listAll)
	mux.Route("/value/{metricType}/{metricName}", func(r chi.Router) {
		r.Use(metricTypeCtx)
		r.Use(c.metricNameCtx)
		r.Get("/", c.get)
	})
	mux.Route("/update/{metricType}/{metricName}/{metricValue}", func(r chi.Router) {
		r.Use(textPlainTypeCheck)
		r.Use(metricTypeCtx)
		r.Use(c.metricNameCtx)
		r.Use(metricValueCtx)
		r.Post("/", c.update)
	})
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseD := &responseData{status: http.StatusOK}
		lw := responseWriter{ResponseWriter: w, responseData: responseD}

		next.ServeHTTP(&lw, r)

		log.Info().
			Str("request_method", r.Method).
			Str("request_uri", r.RequestURI).
			Dur("duration", time.Since(start)).
			Msg("")
		log.Info().
			Int("response_status", responseD.status).
			Int("response_size", responseD.size).
			Msg("")
	})
}

func metricTypeCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, string(metricTypeKey))

		switch metricType {
		case models.Gauge, models.Counter:
			ctx := context.WithValue(r.Context(), metricTypeKey, metricType)
			next.ServeHTTP(w, r.WithContext(ctx))
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})
}

func (c *metricsController) metricNameCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, string(metricNameKey))
		
		if metricName == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var ctx context.Context
		metric, ok := c.storage.Get(r.Context().Value(metricTypeKey).(string) + metricName)

		if ok {
			ctx = context.WithValue(r.Context(), metricKey, metric)
		} else {
			ctx = context.WithValue(r.Context(), metricNameKey, metricName)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func metricValueCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricValue := chi.URLParam(r, string(metricValueKey))

		if metricValue == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), metricValueKey, metricValue)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func textPlainTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func textPlainTypeCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !(strings.HasPrefix(r.Header.Get("Content-Type"), "text/plain") || (r.Header.Get("Content-Type") == "")) {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (c *metricsController) listAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;")
	var body string
	body += "<html><body>"
	for v := range c.storage.All() {
		body += fmt.Sprintf("<p>%s</p>", v.String())
	}
	body += "</body></html>"
	w.Write([]byte(body))
}

func (c *metricsController) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	metric, metricOk := ctx.Value(metricKey).(*models.Metrics)

	if metricOk {
		w.Write([]byte(metric.StringValue()))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (c *metricsController) update(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()
	metric, metricOk := ctx.Value(metricKey).(*models.Metrics)
	newMetric := !metricOk

	if newMetric {
		metric, err = models.NewMetric(ctx.Value(metricTypeKey).(string), ctx.Value(metricNameKey).(string))
		if err == nil {
			metricOk = true
		}
	}

	metricValue, valOk := ctx.Value(metricValueKey).(string)

	err = metric.Update(metricValue)

	if metricOk && valOk && (err == nil) {
		if newMetric { c.storage.Set(metric.MType + metric.ID, metric) }
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
