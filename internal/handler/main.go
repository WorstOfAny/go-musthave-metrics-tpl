package handler

import(
	"net/http"
	"strconv"
)

type MemStorage struct {
	gauges map[string]float64
	counters map[string]int64
}

var storage MemStorage = MemStorage{gauges: make(map[string]float64), counters: make(map[string]int64)}

func update(w http.ResponseWriter, r *http.Request) {
	metric_type := r.PathValue("type")
	metric_name := r.PathValue("varName")
	metric_value := r.PathValue("varValue")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	switch metric_type {
	case "gauge":
		metric_value, err := strconv.ParseFloat(metric_value, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			storage.gauges[metric_name] = metric_value
		}
	case "counter":
		metric_value, err := strconv.ParseInt(metric_value, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			storage.counters[metric_name] += metric_value
		}
	default:
			w.WriteHeader(http.StatusBadRequest)
	}
}

func UpdateHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/{type}/{varName}/{varValue}", update)
	return mux
}
