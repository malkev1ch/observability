package swaggergui

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/vearutop/statigz"

	"github.com/malkev1ch/observability/pkg/swaggergui/static"
)

// Handler handles swagger UI request.
type Handler struct {
	Config

	ConfigJSON template.JS

	tpl          *template.Template
	staticServer http.Handler
}

// NewHandlerWithConfig returns a HTTP handler for swagger UI.
func NewHandlerWithConfig(config Config) *Handler {
	config.BasePath = strings.TrimSuffix(config.BasePath, "/") + "/"

	h := &Handler{
		Config: config,
	}

	j, err := json.Marshal(h.Config)
	if err != nil {
		panic(err)
	}

	h.ConfigJSON = template.JS(j) //nolint:gosec // Data is well-formed.

	h.tpl, err = template.New("index").Parse(IndexTpl(config))
	if err != nil {
		panic(err)
	}

	h.staticServer = http.StripPrefix(h.BasePath, statigz.FileServer(static.FS))

	return h
}

// ServeHTTP implements http.Handler interface to handle swagger UI request.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.TrimSuffix(r.URL.Path, "/") != strings.TrimSuffix(h.BasePath, "/") && h.staticServer != nil {
		h.staticServer.ServeHTTP(w, r)

		return
	}

	w.Header().Set("Content-Type", "text/html")

	if err := h.tpl.Execute(w, h); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
