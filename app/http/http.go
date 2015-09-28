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
	"strings"
	"text/template"
	"time"
	"unicode"

	"github.com/justinas/alice"
	"github.com/supinf/reinvent-sessions-api/app/config"
	"github.com/supinf/reinvent-sessions-api/app/logs"
	"github.com/supinf/reinvent-sessions-api/app/misc"
	"gopkg.in/throttled/throttled.v1"
	"gopkg.in/throttled/throttled.v1/store"
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

// RequestGetParamS retrives a request parameter as string
func RequestGetParamS(r *http.Request, key, def string) string {
	value, found := RequestGetParam(r, key)
	if !found {
		return def
	}
	return value
}

// RequestGetParamI retrives a request parameter as int
func RequestGetParamI(r *http.Request, key string, def int) int {
	value, found := RequestGetParam(r, key)
	if !found {
		return def
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return def
	}
	return i
}

// SplittedUpperStrings split word to array and change those words to UpperCase
func SplittedUpperStrings(value string) []string {
	splitted := strings.FieldsFunc(value, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})
	words := make([]string, len(splitted))
	for i, val := range splitted {
		words[i] = strings.ToUpper(val)
	}
	return words
}

// RequestPostParam retrives a request parameter
func RequestPostParam(r *http.Request, key string) (string, bool) {
	value := r.PostFormValue(key)
	return value, (value != "")
}

// SetCookie set a cookie
func SetCookie(key, value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:   key,
		Value:  value,
		MaxAge: maxAge,
		Secure: cfg.SecuredCookie,
		Path:   "/",
	}
}

// GetCookie retrives a cookie value
func GetCookie(r *http.Request, key string) (id string, err error) {
	var c *http.Cookie
	if c, err = r.Cookie(key); err != nil {
		return
	}
	id = c.Value
	return
}

// Chain enables middleware chaining
func Chain(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return chain(true, false, true, f)
}

// AssetsChain enables middleware chaining
func AssetsChain(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return chain(false, true, false, f)
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
	stage := cfg.Stage
	if !misc.ZeroOrNil(stage) {
		stage = stage + "/"
	}
	if err := tmpl.Execute(w, struct {
		AppName        string
		AppStage       string
		StaticFileHost string
		Data           interface{}
	}{cfg.Name, stage, cfg.StaticFileHost, data}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logs.Error.Printf("ERROR: @RenderHTML %s", err.Error())
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
		logs.Error.Printf("ERROR: %s %s", caption, err.Error())
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

func chain(log, cors, validate bool, f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	chain := alice.New(timeout)
	if th != nil {
		chain = alice.New(th.Throttle, timeout)
	}
	return chain.Then(http.HandlerFunc(custom(log, cors, validate, f)))
}

func custom(log, cors, validate bool, f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		addr := r.RemoteAddr
		if ip, found := header(r, "X-Forwarded-For"); found {
			addr = ip
		}

		// reject if the request is invalid
		if validate {
			if (!misc.ZeroOrNil(cfg.ValidHost) && !strings.Contains(r.Host, cfg.ValidHost)) ||
				(!misc.ZeroOrNil(cfg.ValidUserAgent) && !strings.Contains(r.UserAgent(), cfg.ValidUserAgent)) {
				logs.Info.Printf("%s %s %s %s", addr, strconv.Itoa(http.StatusForbidden), r.Method, r.URL)
				http.Error(w, "403 Forbidden", http.StatusForbidden)
				return
			}
		}
		// compress settings
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
		writer := &customResponseWriter{Writer: ioWriter, ResponseWriter: w, status: 200}

		// CORS headers
		if cors && !misc.ZeroOrNil(cfg.CorsMethods) {
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Methods", cfg.CorsMethods)
			w.Header().Set("Access-Control-Allow-Origin", cfg.CorsOrigin)
		}

		// route to the controllers
		f(writer, r)

		// access log
		if log && cfg.AccessLog {
			logs.Info.Printf("%s %s %s %s", addr, strconv.Itoa(writer.status), r.Method, r.URL)
		}
	}
}

func header(r *http.Request, key string) (string, bool) {
	if r.Header == nil {
		return "", false
	}
	if candidate := r.Header[key]; !misc.ZeroOrNil(candidate) {
		return candidate[0], true
	}
	return "", false
}

func timeout(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, cfg.Timeout*time.Second, "timed out")
}
