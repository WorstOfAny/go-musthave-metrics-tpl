package handler

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/model"
	"net/http"
	"strings"
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
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
