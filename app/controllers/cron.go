package controllers

import (
	"net/http"
	"time"

	"github.com/robfig/cron"
	"github.com/supinf/reinvent-sessions-api/app/crons"
	misc "github.com/supinf/reinvent-sessions-api/app/http"
	"github.com/supinf/reinvent-sessions-api/app/models"
)

func init() {

	http.Handle("/cron/list", misc.Chain(func(w http.ResponseWriter, r *http.Request) {
		entries := []*cron.Entry{}
		for _, c := range crons.Crons() {
			entries = append(entries, c.Entries()...)
		}
		next := []time.Time{}
		for _, entry := range entries {
			next = append(next, entry.Next)
		}
		misc.RenderJSON(w, next, nil)
	}))

	http.Handle("/cron/results", misc.Chain(func(w http.ResponseWriter, r *http.Request) {
		results, _, err := models.GetCronResults()
		misc.RenderJSON(w, results, err)
	}))

}
