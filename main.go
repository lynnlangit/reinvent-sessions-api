package main

import (
	"fmt"
	"net/http"
	"path"

	"github.com/supinf/reinvent-sessions-api/app/config"
	_ "github.com/supinf/reinvent-sessions-api/app/controllers"
	_ "github.com/supinf/reinvent-sessions-api/app/crons"
	misc "github.com/supinf/reinvent-sessions-api/app/http"
	"github.com/supinf/reinvent-sessions-api/app/logs"
	v "github.com/supinf/reinvent-sessions-api/app/misc"
	_ "github.com/supinf/reinvent-sessions-api/app/models"
)

func main() {
	cfg := config.NewConfig()
	logs.Debug("[config] " + cfg.String())

	http.Handle("/", index())
	http.HandleFunc("/alive", alive)
	http.HandleFunc("/version", version)
	http.Handle("/assets/", assets(cfg))

	logs.Infof("[service] listening on port %v", cfg.Port)
	logs.Fatal(http.ListenAndServe(":"+fmt.Sprint(cfg.Port), nil))
}

func index() http.Handler {
	return misc.Chain(true, false, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		misc.RenderHTML(w, []string{"app/index.tmpl"}, nil, nil)
	})
}
func alive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "version: %s", v.Version)
}
func assets(cfg *config.Config) http.Handler {
	fs := http.FileServer(http.Dir(path.Join(cfg.StaticFilePath, "assets")))
	return misc.Chain(false, true, http.StripPrefix("/assets/", fs).ServeHTTP)
}
