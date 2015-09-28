package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kurrik/oauth1a"
	"github.com/supinf/reinvent-sessions-api/app/config"
	util "github.com/supinf/reinvent-sessions-api/app/http"
	"github.com/supinf/reinvent-sessions-api/app/logs"
	"github.com/supinf/reinvent-sessions-api/app/models"
)

var (
	sessions map[string]*oauth1a.UserConfig
	twitter  *oauth1a.Service
)

const (
	twitterTemporaryKey = "tw-temp"
	twitterSessionKey   = "tw-sess"
)

func init() {
	sessions = map[string]*oauth1a.UserConfig{}
	cfg := config.NewConfig()
	twitter = &oauth1a.Service{
		RequestURL:   "https://api.twitter.com/oauth/request_token",
		AuthorizeURL: "https://api.twitter.com/oauth/authorize",
		AccessURL:    "https://api.twitter.com/oauth/access_token",
		ClientConfig: &oauth1a.ClientConfig{
			ConsumerKey:    cfg.TwitterKey,
			ConsumerSecret: cfg.TwitterSecret,
			CallbackURL:    cfg.TwitterCallback,
		},
		Signer: new(oauth1a.HmacSha1Signer),
	}

	/**
	 * Twitter OAuth
	 */
	http.Handle("/twitter/signin", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, util.SetCookie(twitterSessionKey, "", -1))
		_, sessionID, ok := checkTwitterSession(w, r)
		if !ok {
			logs.Error.Print("Could not generate new sessionID.")
			http.Error(w, "Problem generating new session", http.StatusInternalServerError)
			return
		}
		session := &oauth1a.UserConfig{}
		if err := session.GetRequestToken(twitter, new(http.Client)); err != nil {
			logs.Error.Printf("Could not get request token: %v", err)
			http.Error(w, fmt.Sprintf("Problem getting the request token: %v", err), http.StatusInternalServerError)
			return
		}
		url, err := session.GetAuthorizeURL(twitter)
		if err != nil {
			logs.Error.Printf("Could not get authorization URL: %v", err)
			http.Error(w, "Problem getting the authorization URL", http.StatusInternalServerError)
			return
		}
		sessions[sessionID] = session
		http.Redirect(w, r, url, http.StatusFound)
	}))

	http.Handle("/twitter/callback", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		_, found := util.RequestGetParam(r, "denied")
		if found {
			http.Redirect(w, r, "/reinvent/sessions?output=html", http.StatusFound)
			return
		}
		authorizedToken, sessionID, ok := checkTwitterSession(w, r)
		if !ok {
			logs.Error.Print("Could not generate new sessionID.")
			http.Error(w, "Problem generating new session", http.StatusInternalServerError)
			return
		}
		if authorizedToken.ID != "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		session, ok := sessions[sessionID]
		if !ok {
			logs.Error.Print("Could not find user config in sesions storage.")
			http.Error(w, "Invalid session", http.StatusBadRequest)
			return
		}
		token, verifier, err := session.ParseAuthorize(r, twitter)
		if err != nil {
			logs.Error.Printf("Could not parse authorization: %v", err)
			http.Error(w, "Problem parsing authorization", http.StatusInternalServerError)
			return
		}
		if err = session.GetAccessToken(token, verifier, twitter, new(http.Client)); err != nil {
			logs.Error.Printf("Error getting access token: %v", err)
			http.Error(w, "Problem getting an access token", http.StatusInternalServerError)
			return
		}
		delete(sessions, sessionID)

		authorized := models.OAuthAuthorizedToken{
			ID:                session.AccessValues.Get("user_id"),
			ScreenName:        session.AccessValues.Get("screen_name"),
			AccessTokenKey:    session.AccessTokenKey,
			AccessTokenSecret: session.AccessTokenSecret,
			CognitoPoolID:     strings.Replace(cfg.CognitoPoolID, ":", "*", -1),
			CognitoRoleArn:    strings.Replace(cfg.CognitoRoleArn, ":", "*", -1),
		}
		bytes, _ := json.Marshal(authorized)
		http.SetCookie(w, util.SetCookie(twitterSessionKey, string(bytes), 60*60*24))
		http.SetCookie(w, util.SetCookie(twitterTemporaryKey, "", -1))
		http.Redirect(w, r, "/reinvent/sessions?output=html", http.StatusFound)
	}))
}

func checkTwitterSession(w http.ResponseWriter, r *http.Request) (token models.OAuthAuthorizedToken, tempID string, found bool) {
	if value, err := util.GetCookie(r, twitterSessionKey); err == nil {
		value = strings.Replace(strings.Replace(strings.Replace(value, ":", "\":\"", -1), ",", "\",\"", -1), "*", ":", -1)
		value = strings.Replace(strings.Replace(value, "{", "{\"", -1), "}", "\"}", -1)
		if err = json.Unmarshal([]byte(value), &token); err != nil {
			logs.Error.Printf("Error: %v", err)
			return
		}
		return token, "", true
	}
	if value, err := util.GetCookie(r, twitterTemporaryKey); err == nil {
		logs.Info.Printf("temp: %v", value)
		return token, value, true
	}
	sessionID, ok := newTwitterSessionID()
	if !ok {
		logs.Error.Print("Could not generate new sessionID.")
		return token, "", false
	}
	logs.Error.Print("!!?.")
	http.SetCookie(w, util.SetCookie(twitterTemporaryKey, sessionID, 60*60*24))
	return token, "", false
}

func newTwitterSessionID() (string, bool) {
	b := make([]byte, 128)
	n, err := io.ReadFull(rand.Reader, b)
	if n != len(b) || err != nil {
		return "", false
	}
	return base64.URLEncoding.EncodeToString(b), true
}
