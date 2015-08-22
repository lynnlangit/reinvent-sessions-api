package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

	misc "github.com/supinf/reinvent-sessions-api/app/http"
	"github.com/supinf/reinvent-sessions-api/app/models"
)

func init() {

	/**
	 * sessions
	 * @param string output [The formatting style for response body (html | json).]
	 * @param string q [Space seperated words to use in filtering the response data (for example, best practice).]
	 */
	http.Handle("/reinvent-sessions", misc.Chain(func(w http.ResponseWriter, r *http.Request) {
		output, found := misc.RequestGetParam(r, "output")
		if found && (strings.ToLower(output) == "html") {
			misc.RenderHTML(w, []string{"reinvent/index.tmpl"}, nil, nil)
			return
		}
		cron, err1 := models.GetCronResult("SyncReInventSessions")
		if misc.IsInvalid(w, err1, "@aws.DynamoRecord") {
			return
		}
		// get sessions
		sessions, _, err2 := models.GetSessions()
		if misc.IsInvalid(w, err2, "@aws.DynamoRecords") {
			return
		}
		// filter
		if q, found := misc.RequestGetParam(r, "q"); found {
			splitted := strings.FieldsFunc(q, func(c rune) bool {
				return !unicode.IsLetter(c) && !unicode.IsNumber(c)
			})
			words := make([]string, len(splitted))
			for i, val := range splitted {
				words[i] = strings.ToUpper(val)
			}
			filtered := models.Sessions{}
			for _, session := range sessions {
				if session.Contains(words) {
					filtered = append(filtered, session)
				}
			}
			sessions = filtered
		}
		misc.RenderJSON(w, struct {
			Count    int              `json:"count"`
			Sessions []models.Session `json:"sessions"`
			Sync     time.Time        `json:"sync"`
		}{
			Count:    len(sessions),
			Sessions: sessions,
			Sync:     cron.LastEndDate,
		}, nil)
	}))

	/**
	 * session
	 * @param string id [Session ID]
	 */
	http.Handle("/reinvent-session", misc.Chain(func(w http.ResponseWriter, r *http.Request) {
		id, found := misc.RequestGetParam(r, "id")
		if !found {
			fmt.Print(w, "Parameter [ id ] is needed.")
			return
		}
		session, err := models.GetSession(id)
		if misc.IsInvalid(w, err, "@aws.DynamoRecords") {
			return
		}
		misc.RenderJSON(w, session, nil)
	}))
}
