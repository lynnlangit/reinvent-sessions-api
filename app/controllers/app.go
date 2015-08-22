// Package controllers implements functions to route user requests
package controllers

import (
	"net/http"

	misc "github.com/supinf/reinvent-sessions-api/app/http"
	"github.com/supinf/reinvent-sessions-api/app/models"
)

func init() {

	http.Handle("/api-list", misc.Chain(func(w http.ResponseWriter, r *http.Request) {
		misc.RenderJSON(w, models.ListAPI(), nil)
	}))

}
