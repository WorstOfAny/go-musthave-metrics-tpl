package handler

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"context"
	"fmt"
	"iter"
	"time"
	"encoding/json"
	"net/http/httputil"
	"compress/gzip"
	"io"
)

type metricsStorage interface {
 Set(string, *models.Metrics)
 Get(string) (*models.Metrics, bool)
 Remove(string)
 All() iter.Seq[*models.Metrics]
}

type DB interface {
	Ping() error
}

type metricsController struct {
	storage metricsStorage
	db DB
}

func NewMetricsController(storage metricsStorage, db DB) *metricsController {
	return &metricsController{storage: storage, db: db}
}

type responseData struct {
	status int
	size int
}
type responseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
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
	mux.Use(decodeRequest)
	mux.Use(logRequest)
	mux.With(textPlainTypeCheck, encodeResponse).Get("/", c.listAll)
	mux.With(textPlainTypeSet).Get("/ping", c.ping)
	mux.Route("/value", func(r chi.Router) {
		r.With(jsonTypeSet, jsonTypeCheck, c.jsonCtx, encodeResponse).Post("/", c.showJSON)
		r.Route("/{metricType}/{metricName}", func(r chi.Router){
			r.Use(textPlainTypeSet, textPlainTypeCheck, metricTypeCtx, c.metricNameCtx)
			r.Get("/", c.showTextPlain)
		})
	})

	mux.Route("/update", func(r chi.Router) {
		r.With(jsonTypeSet, jsonTypeCheck, c.jsonCtx, encodeResponse).Post("/", c.update)
		r.Route("/{metricType}/{metricName}/{metricValue}", func(r chi.Router) {
			r.Use(textPlainTypeSet, textPlainTypeCheck, metricTypeCtx, c.metricNameCtx, metricValueCtx)
			r.Post("/", c.update)
		})
	})
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqDump, err := httputil.DumpRequest(r, true)

		if err != nil {
			log.Debug().
				Err(err).
				Msg("dump request error")
		}

		fmt.Println(string(reqDump))
		responseD := &responseData{status: http.StatusOK}
		lw := responseWriter{ResponseWriter: w, responseData: responseD}

		next.ServeHTTP(&lw, r)

		log.Info().
			Str("request_method", r.Method).
			Str("request_uri", r.RequestURI).
			Str("request_content_type", r.Header.Get("Content-Type")).
			Dur("duration", time.Since(start)).
			Int("response_status", responseD.status).
			Int("response_size", responseD.size).
			Str("response_content_type", lw.Header().Get("Content-Type")).
			Msg("")
	})
}

func decodeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)

		if err != nil {
			log.Debug().
				Err(err).
				Msg("gzip reader error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer gz.Close()

		r.Body = gz
		next.ServeHTTP(w, r)
	})
}

func encodeResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)

		if err != nil {
			log.Debug().
				Err(err).
				Msg("gzip writer error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func (c *metricsController) jsonCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqMetric models.Metrics
		dec := json.NewDecoder(r.Body)

		if err := dec.Decode(&reqMetric); err != nil {
			log.Debug().
				Err(err).
				Msg("json decode error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch reqMetric.MType {
		case models.Counter, models.Gauge:
			ctx := context.WithValue(r.Context(), metricTypeKey, reqMetric.MType)
			ctx = context.WithValue(ctx, metricNameKey, reqMetric.ID)
			ctx = context.WithValue(ctx, metricValueKey, reqMetric.StringValue())
			metric, found := c.storage.Get(reqMetric.MType + reqMetric.ID)

			if found {
				ctx = context.WithValue(ctx, metricKey, metric)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
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
			w.Header().Set("Content-Type", "text/plain;")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func jsonTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;")
		next.ServeHTTP(w, r)
	})
}

func jsonTypeCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !(strings.HasPrefix(r.Header.Get("Content-Type"), "application/json")) {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (c *metricsController) listAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var body string
	body += "<html><body>"
	for v := range c.storage.All() {
		body += fmt.Sprintf("<p>%s</p>", v.String())
	}
	body += "</body></html>"
	w.Write([]byte(body))
}

func (c *metricsController) showJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	metric, metricOk := ctx.Value(metricKey).(*models.Metrics)

	if metricOk {
		body, err := json.Marshal(metric)

		if err != nil {
			log.Debug().
				Err(err).
				Str("metricID", metric.ID).
				Str("metricType", metric.MType).
				Msg("marshal error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(body)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
	
}

func (c *metricsController) showTextPlain(w http.ResponseWriter, r *http.Request) {
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
		response, err := json.Marshal(metric)

		if err != nil {
			log.Debug().
				Err(err).
				Str("metricID", metric.ID).
				Str("metricType", metric.MType).
				Msg("marshal error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(response)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (c *metricsController) ping(w http.ResponseWriter, r *http.Request) {
	if err := c.db.Ping(); err != nil {
			log.Debug().
				Err(err).
				Msg("failed to ping database")
			w.WriteHeader(http.StatusInternalServerError)
			return
	}

	w.Write([]byte("Success"))
}
