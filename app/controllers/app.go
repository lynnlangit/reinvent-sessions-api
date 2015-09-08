// Package controllers implements functions to route user requests
package controllers

import (
	"net/http"

	util "github.com/supinf/reinvent-sessions-api/app/http"
	"github.com/supinf/reinvent-sessions-api/app/models"
)

func init() {

	http.Handle("/api-list", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		util.RenderJSON(w, models.ListAPI(), nil)
	}))

}
