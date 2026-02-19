package handler

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strings"
	"context"
	"fmt"
)

type metricsController struct {}

func (c *metricsController) ApplyTo(mux chi.Router) {
	mux.Use(middleware.Logger)
	mux.Use(textPlainTypeSet)
	mux.Get("/", c.listAll)
	mux.Route("/value/{metricType}/{metricName}", func(r chi.Router) {
		r.Use(metricTypeCtx)
		r.Use(metricNameCtx)
		r.Get("/", c.get)
	})
	mux.Route("/update/{metricType}/{metricName}/{metricValue}", func(r chi.Router) {
		r.Use(textPlainTypeCheck)
		r.Use(metricTypeCtx)
		r.Use(metricNameCtx)
		r.Use(metricValueCtx)
		r.Post("/", c.update)
	})
}

func metricTypeCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")

		switch metricType {
		case models.Gauge, models.Counter:
			ctx := context.WithValue(r.Context(), "metricType", metricType)
			next.ServeHTTP(w, r.WithContext(ctx))
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})
}

func metricNameCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, "metricName")
		
		if metricName == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var ctx context.Context
		metric, ok := models.FindMetric(r.Context().Value("metricType").(string), metricName)

		if ok {
			ctx = context.WithValue(r.Context(), "metric", metric)
		} else {
			ctx = context.WithValue(r.Context(), "metricName", metricName)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func metricValueCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricValue := chi.URLParam(r, "metricValue")

		if metricValue == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "metricValue", metricValue)
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

func NewMetricsController() metricsController {
	return metricsController{}
}

func (c *metricsController) listAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;")
	var body string
	body += "<html><body>"
	for v := range models.AllMetrics() {
		body += fmt.Sprintf("<p>%s</p>", v.String())
	}
	body += "</body></html>"
	w.Write([]byte(body))
}

func (c *metricsController) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	metric, metricOk := ctx.Value("metric").(*models.Metrics)

	if metricOk {
		w.Write([]byte(metric.StringValue()))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (c *metricsController) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	metric, metricOk := ctx.Value("metric").(*models.Metrics)
	newMetric := !metricOk

	if !metricOk {
		metric, metricOk = models.NewMetric(ctx.Value("metricType").(string), ctx.Value("metricName").(string))
	}

	metricValue, valOk := ctx.Value("metricValue").(string)

	if metricOk && valOk && metric.Update(metricValue) {
		if newMetric { metric.Save() }
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
