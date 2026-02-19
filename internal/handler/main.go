package handler

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/logger"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/model"
	"net/http"
	"strings"
)

type controller struct {
	l logger.Logger
}

func (c *controller) ApplyTo(mux *http.ServeMux) {
	mux.Handle("/update/{type}/{varName}/{varValue}", http.HandlerFunc(c.Update))
}

func NewController(l logger.Logger) controller {
	return controller{l: l}
}

func (c *controller) Update(w http.ResponseWriter, r *http.Request) {
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
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
