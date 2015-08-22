// Package http provides a convenient way to impliment http servers
package http

import (
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"html"
	"io"
	"net/http"
	"path"
	"strconv"
	"text/template"
	"time"

	"github.com/justinas/alice"
	"github.com/supinf/reinvent-sessions-api/app/config"
	"github.com/supinf/reinvent-sessions-api/app/logs"
	"github.com/supinf/reinvent-sessions-api/app/misc"
	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store"
)

var cfg *config.Config
var th *throttled.Throttler

func init() {
	cfg = config.NewConfig()

	// Access limitations
	if cfg.LimitRatePerMin > 0 {
		if cfg.LimitBursts > 0 {
			th = throttled.Interval(
				throttled.PerMin(cfg.LimitRatePerMin),
				cfg.LimitBursts, cfg.LimitVaryBy,
				cfg.LimitKeyCache)
		} else {
			th = throttled.RateLimit(
				throttled.PerMin(cfg.LimitRatePerMin),
				cfg.LimitVaryBy, store.NewMemStore(cfg.LimitKeyCache))
		}
	}
}

// RequestGetParam retrives a request parameter
func RequestGetParam(r *http.Request, key string) (string, bool) {
	value := r.URL.Query().Get(key)
	return value, (value != "")
}

// RequestPostParam retrives a request parameter
func RequestPostParam(r *http.Request, key string) (string, bool) {
	value := r.PostFormValue(key)
	return value, (value != "")
}

// Chain enables middleware chaining
func Chain(log, cors bool, f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	chain := alice.New(timeout)
	if th != nil {
		chain = alice.New(th.Throttle, timeout)
	}
	return chain.Then(http.HandlerFunc(custom(f, log, cors)))
}

// RenderText write data as a simple text
func RenderText(w http.ResponseWriter, data string, err error) {
	if IsInvalid(w, err, "@RenderText") {
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(html.EscapeString(data)))
}

// RenderHTML write data as a HTML text with template
func RenderHTML(w http.ResponseWriter, templatePath []string, data interface{}, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	relatives := append([]string{"base.tmpl"}, templatePath...)
	templates := make([]string, len(relatives))
	for idx, template := range relatives {
		templates[idx] = path.Join(cfg.StaticFilePath, "views", template)
	}

	tmpl, err := template.ParseFiles(templates...)
	if IsInvalid(w, err, "@RenderHTML") {
		return
	}
	// if err := tmpl.Execute(w, data, cfg.StaticFileHost); err != nil {
	if err := tmpl.Execute(w, struct {
		AppName        string
		StaticFileHost string
		Data           interface{}
	}{cfg.Name, cfg.StaticFileHost, data}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logs.Errorf("ERROR: @RenderHTML %s", err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/html")
}

// RenderJSON write data as a json
func RenderJSON(w http.ResponseWriter, data interface{}, err error) {
	if IsInvalid(w, err, "@RenderJSON") {
		return
	}
	js, err := json.Marshal(data)
	if IsInvalid(w, err, "@RenderJSON") {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// IsInvalid checks if the second argument represents a real error
func IsInvalid(w http.ResponseWriter, err error, caption string) bool {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logs.Errorf("ERROR: %s %s", caption, err.Error())
		return true
	}
	return false
}

type customResponseWriter struct {
	io.Writer
	http.ResponseWriter
	status int
}

func (r *customResponseWriter) Write(b []byte) (int, error) {
	if r.Header().Get("Content-Type") == "" {
		r.Header().Set("Content-Type", http.DetectContentType(b))
	}
	return r.Writer.Write(b)
}

func (r *customResponseWriter) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.status = status
}

func custom(f func(w http.ResponseWriter, r *http.Request), log, cors bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ioWriter := w.(io.Writer)
		for _, val := range misc.ParseCsvLine(r.Header.Get("Accept-Encoding")) {
			if val == "gzip" {
				w.Header().Set("Content-Encoding", "gzip")
				g := gzip.NewWriter(w)
				defer g.Close()
				ioWriter = g
				break
			}
			if val == "deflate" {
				w.Header().Set("Content-Encoding", "deflate")
				z := zlib.NewWriter(w)
				defer z.Close()
				ioWriter = z
				break
			}
		}
		if !misc.ZeroOrNil(cfg.CorsMethods) {
			w.Header().Set("Access-Control-Allow-Headers", "'*'")
			w.Header().Set("Access-Control-Allow-Methods", "'"+cfg.CorsMethods+"'")
			w.Header().Set("Access-Control-Allow-Origin", "''"+cfg.CorsOrigin+"'")
		}
		writer := &customResponseWriter{Writer: ioWriter, ResponseWriter: w, status: 200}
		f(writer, r)
		if log && cfg.AccessLog {
			logs.Infof("%s %s %s %s", r.RemoteAddr, strconv.Itoa(writer.status), r.Method, r.URL)
		}
	}
}

func timeout(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, cfg.Timeout*time.Second, "timed out")
}
