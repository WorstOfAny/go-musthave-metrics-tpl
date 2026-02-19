package handler

import(
	"internal/logger"
	"net/http"
	"strings"
)
type Setter interface { Set(name string, value string) error }
type Getter interface { Get(name string) (any, bool) }
type Remover interface { Remove(name string) }
type Repository interface {
	Setter
	Getter
	Remover
}

type controller struct {
	l logger.Logger
	repositories map[string]Repository
}

func (c *controller) ApplyTo(mux *http.ServeMux) {
	mux.Handle("/update/{type}/{varName}/{varValue}", http.HandlerFunc(c.Update))
}

func NewController(l logger.Logger, repositories map[string]Repository) controller {
	return controller{l: l, repositories: repositories}
}

func (c *controller) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "text/plain") {
		http.Error(w, "", http.StatusUnsupportedMediaType)
		return
	}

	metric_type := r.PathValue("type")
	metric_name := r.PathValue("varName")
	metric_value := r.PathValue("varValue")

	c.l.Log(r.URL)
	c.l.Log(metric_type)
	c.l.Log(metric_name)
	c.l.Log(metric_value)

	if repo, ok := c.repositories[metric_type]; ok {
		 err := repo.Set(metric_name, metric_value)
		 if err != nil {
				w.WriteHeader(http.StatusBadRequest)
		 }
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
