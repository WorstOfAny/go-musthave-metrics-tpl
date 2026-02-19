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

type metricsController struct {
}

func (c *metricsController) ApplyTo(mux *http.ServeMux) {
	mux.Handle("/update/{type}/{varName}/{varValue}", http.HandlerFunc(c.Update))
}

func NewController() *metricsController {
	return &metricsController{}
}

func (c *metricsController) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "text/plain") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	metric, ok := models.FindMetric(r.PathValue("type"), r.PathValue("varName"))
	newMetric := !ok

	if !ok {
		metric, ok = models.NewMetric(r.PathValue("type"), r.PathValue("varName"))
	}

	if ok && metric.Update(r.PathValue("varValue")) {
		if newMetric {
			metric.Save()
		}
	})
}

func metricNameCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, string(metricNameKey))
		
		if metricName == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var ctx context.Context
		metric, ok := models.FindMetric(r.Context().Value(metricTypeKey).(string), metricName)

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
	metric, metricOk := ctx.Value(metricKey).(*models.Metrics)

	if metricOk {
		w.Write([]byte(metric.StringValue()))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (c *metricsController) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	metric, metricOk := ctx.Value(metricKey).(*models.Metrics)
	newMetric := !metricOk

	if !metricOk {
		metric, metricOk = models.NewMetric(ctx.Value(metricTypeKey).(string), ctx.Value(metricNameKey).(string))
	}

	metricValue, valOk := ctx.Value(metricValueKey).(string)

	if metricOk && valOk && metric.Update(metricValue) {
		if newMetric { metric.Save() }
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
