package handler

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/http/httputil"
	"net/http/pprof"
	"runtime/debug"
	"slices"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	pgconn "github.com/jackc/pgconn"
	"github.com/rs/zerolog/log"

	models "github.com/WorstOfAny/go-musthave-metrics-tpl/internal/model"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/observers"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/repository"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/service"
)

type metricsController struct {
	storage    repository.Repository
	observers  []observers.Observer
	privateKey *rsa.PrivateKey
	key        string
}

// NewMetricsController конструктор для metricsController, возвращает указатель на metricsController
func NewMetricsController(ctx context.Context, storage repository.Repository, key string, privateBytes []byte) *metricsController {
	controller := &metricsController{storage: storage, key: key}
	controller.observers = make([]observers.Observer, 0)

	privatePemBlock, _ := pem.Decode(privateBytes)
	if privatePemBlock == nil {
		log.Error().Msg("private key not found")
		return controller
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privatePemBlock.Bytes)
	if err != nil {
		log.Error().Err(err).Msg("failed parse private key")
		return controller
	}
	controller.privateKey = privateKey
	return controller
}

// Event событие, которое наблюдатель отслеживает
type Event struct {
	Metrics []string
	TS      time.Time
	IPAddr  string
}

// Register добавление наблюдателей для отслеживания событий
func (c *metricsController) Register(o observers.Observer) {
	c.observers = append(c.observers, o)
}

func (c *metricsController) newEvent(metrics []models.Metrics, ip_addr string) {
	metricIDs := make([]string, len(metrics))
	for i, m := range metrics {
		metricIDs[i] = m.ID
	}

	event := Event{Metrics: metricIDs, TS: time.Now(), IPAddr: ip_addr}
	data, err := json.Marshal(event)

	if err != nil {
		log.Error().Err(err).Msg("marshal event error")
		return
	}

	for _, obs := range c.observers {
		obs.Update(data)
	}

}

type responseData struct {
	status int
	size   int
}
type responseWriter struct {
	http.ResponseWriter
	responseData *responseData
	key          string
}

func (wr *responseWriter) Write(b []byte) (int, error) {
	if wr.key != "" {
		signedBody := service.SignToString(b, wr.key)
		wr.ResponseWriter.Header().Set("HashSHA256", signedBody)
	}
	size, err := wr.ResponseWriter.Write(b)
	wr.responseData.size += size
	return size, err
}

func (wr *responseWriter) WriteHeader(statusCode int) {
	wr.ResponseWriter.WriteHeader(statusCode)
	wr.responseData.status = statusCode
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

type ctxKey string

const (
	metricKey  ctxKey = "metric"
	metricsKey ctxKey = "metrics"
)

// ApplyTo настройка маршрутизатора mux для работы с эндпоинтами metricsController
func (c *metricsController) ApplyTo(mux chi.Router) {
	mux.Use(recoveryPanic, c.checkHMAC, c.decryptRequest, decodeRequest, c.logRequest)
	mux.With(textPlainTypeCheck, encodeResponse).Get("/", c.listAll)
	mux.With(textPlainTypeSet).Get("/ping", c.ping)
	mux.Route("/debug/", func(r chi.Router) {
		r.HandleFunc("/pprof", pprof.Index)
		r.HandleFunc("/cmdline", pprof.Cmdline)
		r.HandleFunc("/profile", pprof.Profile)
		r.HandleFunc("/symbol", pprof.Symbol)
		r.HandleFunc("/trace", pprof.Trace)
		r.HandleFunc("/*", http.HandlerFunc(pprof.Index))
	})

	mux.Route("/value", func(r chi.Router) {
		r.With(jsonTypeSet, jsonTypeCheck, c.jsonCtx, encodeResponse).Post("/", c.showJSON)
		r.Route("/{metricType}/{metricName}", func(r chi.Router) {
			r.With(textPlainTypeSet, textPlainTypeCheck, c.plainGetCtx).Get("/", c.showTextPlain)
		})
	})

	mux.Route("/update", func(r chi.Router) {
		r.With(jsonTypeSet, jsonTypeCheck, c.jsonCtx, encodeResponse).Post("/", c.update)
		r.Route("/{metricType}/{metricName}/{metricValue}", func(r chi.Router) {
			r.With(textPlainTypeSet, textPlainTypeCheck, c.plainPostCtx).Post("/", c.update)
		})
	})
	mux.With(jsonTypeSet, jsonTypeCheck, c.bulkJSONCtx, encodeResponse).Post("/updates/", c.updates)
	mux.With(jsonTypeSet, jsonTypeCheck, c.bulkJSONCtx, encodeResponse).Post("/updates", c.updates)
}

func recoveryPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					log.Debug().Err(err).Str("stack", string(debug.Stack())).Msg("panic")
				} else {
					log.Debug().Any("recovered", r).Msg("panic")
				}
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (c *metricsController) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqDump, err := httputil.DumpRequest(r, true)

		if err != nil {
			log.Debug().
				Err(err).
				Msg("dump request error")
		}

		log.Info().Str("req_dump", string(reqDump)).Msg("")
		responseD := &responseData{status: http.StatusOK}
		lw := responseWriter{ResponseWriter: w, responseData: responseD, key: c.key}

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

func (c *metricsController) checkHMAC(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := r.Header.Get("HashSHA256")
		if c.key != "" && msg != "" {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Debug().Err(err).Msg("failed read body")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body.Close()

			err = service.Equal(msg, bodyBytes, c.key)
			if err != nil {
				log.Debug().Err(err).Msg("hmac check fail")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		next.ServeHTTP(w, r)
	})
}

func (c *metricsController) decryptRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c.privateKey == nil {
			next.ServeHTTP(w, r)
			return
		}
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Debug().Err(err).Msg("failed read body")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Body.Close()
		decryptedBody, err := rsa.DecryptPKCS1v15(rand.Reader, c.privateKey, bodyBytes)
		if err != nil {
			log.Debug().Err(err).Msg("failed decrypt body")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(decryptedBody))
		next.ServeHTTP(w, r)
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
			log.Debug().Err(err).Msg("gzip reader error")
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
			log.Debug().Err(err).Msg("gzip writer error")
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
		if !decodeMetrics(w, r, &reqMetric) {
			return
		}

		ctx := context.WithValue(r.Context(), metricKey, reqMetric)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *metricsController) bulkJSONCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqMetrics []models.Metrics
		if !decodeMetrics(w, r, &reqMetrics) {
			return
		}

		ctx := context.WithValue(r.Context(), metricsKey, reqMetrics)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func decodeMetrics(w http.ResponseWriter, r *http.Request, ptr any) bool {
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(ptr); err != nil {
		var metrErr *models.MetricError

		if errors.As(err, &metrErr) {
			if strings.HasPrefix(r.URL.Path, "/value") {
				return true
			} else {
				log.Debug().Err(err).Msg("invalid metric parameters")
			}
		} else {
			log.Debug().Err(err).Msg("json decode error")
		}
		w.WriteHeader(http.StatusBadRequest)

		return false
	}

	return true
}

func (c *metricsController) plainPostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		metricValue := chi.URLParam(r, "metricValue")
		metric, err := models.NewMetric(metricType, metricName, metricValue)
		if err != nil {
			metricErrorHandler(w, r, err, "can't create metric from parameters")
			return
		}

		if ok, err := metric.Valid(); !ok {
			metricErrorHandler(w, r, err, "invalid metric parameters")
			return
		}

		ctx := context.WithValue(r.Context(), metricKey, *metric)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *metricsController) plainGetCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		metric, err := c.storage.Get(r.Context(), metricType+metricName)

		if err != nil {
			metricErrorHandler(w, r, err, "error while fetching metric")
			return
		}

		ctx := context.WithValue(r.Context(), metricKey, metric)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func metricErrorHandler(w http.ResponseWriter, r *http.Request, err error, logMsg string) {
	log.Debug().Err(err).Msg(logMsg)

	if errors.Is(err, repository.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var metrErr *models.MetricError
	if errors.As(err, &metrErr) {
		switch metrErr.Cause() {
		case models.EmptyID:
			w.WriteHeader(http.StatusNotFound)
		case models.EmptyValue, models.WrongType, models.InvalidFloat, models.InvalidInt:
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
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
	var body bytes.Buffer
	body.WriteString("<html><body>")
	objs, err := c.storage.All(r.Context())
	if err != nil {
		log.Debug().Err(err).Msg("failed to fetch metrics from db")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for v := range objs {
		fmt.Fprintf(&body, "<p>%s</p>", v.String())
	}
	body.WriteString("</body></html>")
	w.Write(body.Bytes())
}

func (c *metricsController) showJSON(w http.ResponseWriter, r *http.Request) {
	m, ok := r.Context().Value(metricKey).(models.Metrics)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	metric, err := c.storage.Get(r.Context(), m.Key())

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	body, err := json.Marshal(metric)

	if err != nil {
		log.Debug().Err(err).Str("metricID", metric.ID).Str("metricType", metric.MType).Msg("marshal error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

func (c *metricsController) showTextPlain(w http.ResponseWriter, r *http.Request) {
	m, ok := r.Context().Value(metricKey).(models.Metrics)

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write([]byte(m.StringValue()))
}

func (c *metricsController) update(w http.ResponseWriter, r *http.Request) {
	metric, ok := r.Context().Value(metricKey).(models.Metrics)

	if !ok {
		log.Debug().Msg("Failed to fetch metric from request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if m, err := c.storage.Get(r.Context(), metric.Key()); err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			log.Debug().Err(err).Msg("Failed to fetch metric")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		err := m.Update(metric.StringValue())
		if err != nil {
			log.Debug().Err(err).Msg("Failed to update metric")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metric = m
	}

	err := c.storage.Set(r.Context(), metric)
	if err != nil {
		log.Debug().Err(err).Str("metricID", metric.ID).Str("metricType", metric.MType).Msg("storage set err")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(metric)

	if err != nil {
		log.Debug().Err(err).Str("metricID", metric.ID).Str("metricType", metric.MType).Msg("marshal error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
	go c.newEvent([]models.Metrics{metric}, r.RemoteAddr)
}

func (c *metricsController) updates(w http.ResponseWriter, r *http.Request) {
	metrics, ok := r.Context().Value(metricsKey).([]models.Metrics)
	if !ok {
		log.Debug().Msg("update error: recieved nothing from body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tmpstore := map[string]models.Metrics{}
	for _, metric := range metrics {
		if m, err := c.storage.Get(r.Context(), metric.Key()); err != nil {
			if !errors.Is(err, repository.ErrNotFound) {
				log.Debug().Err(err).Msg("Failed to fetch metric")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			tmpmetric, ok := tmpstore[metric.Key()]
			if ok {
				err := tmpmetric.Update(metric.StringValue())
				if err != nil {
					log.Debug().Err(err).Msg("Failed to update metric")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				tmpstore[metric.Key()] = tmpmetric
			} else {

				tmpstore[metric.Key()] = metric
			}
		} else {
			err := m.Update(metric.StringValue())
			if err != nil {
				log.Debug().Err(err).Msg("Failed to update metric")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			tmpstore[m.Key()] = m
		}
	}

	err := c.storage.BulkSet(r.Context(), slices.Collect(maps.Values(tmpstore)))

	if err != nil {
		log.Debug().Err(err).Msg("storage bulkset err")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(metrics)

	if err != nil {
		log.Debug().Err(err).Msg("marshal error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
	go c.newEvent(metrics, r.RemoteAddr)
}

func (c *metricsController) ping(w http.ResponseWriter, r *http.Request) {
	if err := c.storage.Ping(r.Context()); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Debug().
				Str("pg message", pgErr.Message).
				Err(err).
				Msg("failed to ping database")
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			log.Debug().
				Err(err).
				Msg("failed to ping database")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Write([]byte("Success"))
}
