package models

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/supinf/reinvent-sessions-api/app/logs"

	"gopkg.in/xmlpath.v2"
)

// Path caches xmlpath.Path
type Path struct {
	session       *xmlpath.Path
	sAbbreviation *xmlpath.Path
	sTitle        *xmlpath.Path
	sReference    *xmlpath.Path
	sAbstract     *xmlpath.Path
	sLength       *xmlpath.Path
	sType         *xmlpath.Path
	dTitle        *xmlpath.Path
	dAbstract     *xmlpath.Path
	dTypeID       *xmlpath.Path
	dType         *xmlpath.Path
	dTrackID      *xmlpath.Path
	dTrack        *xmlpath.Path
	dLevelID      *xmlpath.Path
	dLevel        *xmlpath.Path
}

var path Path

func init() {
	path = Path{
		xmlpath.MustCompile(`//div[@class='detailColumn']`),
		xmlpath.MustCompile(`.//span[@class='abbreviation']`),
		xmlpath.MustCompile(`.//span[@class='title']`),
		xmlpath.MustCompile(`./a[@class='openInPopup']/@href`),
		xmlpath.MustCompile(`./span[@class='abstract']`),
		xmlpath.MustCompile(`./small[@class='length']`),
		xmlpath.MustCompile(`./small[@class='type']`),
		xmlpath.MustCompile(`.//div[@class='detailHeader']/h1`),
		xmlpath.MustCompile(`.//fieldset[@id='abstract']/p[2]`),
		xmlpath.MustCompile(`.//div[@id='profileItem_1036_tr']/div[@class='paragraph']/input/@value`),
		xmlpath.MustCompile(`.//div[@id='profileItem_1036_tr']/div[@class='paragraph']`),
		xmlpath.MustCompile(`.//div[@id='profileItem_10042_tr']/div[@class='paragraph']/input/@value`),
		xmlpath.MustCompile(`.//div[@id='profileItem_10042_tr']/div[@class='paragraph']`),
		xmlpath.MustCompile(`.//div[@id='profileItem_10041_tr']/div[@class='paragraph']/input/@value`),
		xmlpath.MustCompile(`.//div[@id='profileItem_10041_tr']/div[@class='paragraph']`),
	}
}

// SyncReInventSessions collect session at AWS re:Invent
func SyncReInventSessions(persist bool) (sessions []Session, err error) {
	logs.Debug("[proc] SyncReInventSessions start")
	requests := []url.Values{}
	start := time.Now()

	// Breakout sessions
	for _, sessionTrack := range sessionTracks {
		data := url.Values{}
		data["key"] = []string{sessionTrack.Caption}
		data["searchType"] = []string{"session"}
		data["value(profileItem_10042)"] = []string{fmt.Sprint(sessionTrack.ID)}
		requests = append(requests, data)
	}
	// The other sessions
	for _, sessionType := range sessionTypes {
		if sessionType.SearchID == 2 {
			continue
		}
		data := url.Values{}
		data["key"] = []string{sessionType.Caption}
		data["searchType"] = []string{"session"}
		data["sessionTypeID"] = []string{fmt.Sprint(sessionType.SearchID)}
		requests = append(requests, data)
	}

	// Request sessions
	for _, request := range requests {
		list, err := getSessions(request)
		if err != nil {
			logs.Errorf("@models.getSessions Error: %s", err)
			PersistCronResult("SyncReInventSessions",
				fmt.Sprintf("length: %v, error: %v", len(sessions), err),
				start, time.Now())
			logs.Debug("[proc] SyncReInventSessions end")
			return sessions, err
		}
		logs.Debugf("[proc] SyncReInventSessions got %v: %v", request["key"], len(list))

		for _, session := range list {
			if sess, err := getSessionDetails(session); err != nil {
				logs.Errorf("@reinvent.getSessionDetails Error: %s", err)
			} else {
				sessions = append(sessions, sess)
			}
		}
	}
	if persist {
		for _, session := range sessions {
			err := session.persist()
			if err != nil {
				logs.Errorf("@models.PersistSession Error: %s", err)
			}
		}
		logs.Debugf("[proc] SyncReInventSessions persisted: %v", len(sessions))
	}
	PersistCronResult("SyncReInventSessions",
		fmt.Sprintf("length: %v, error: %v", len(sessions), "-"),
		start, time.Now())
	logs.Debug("[proc] SyncReInventSessions end")
	return sessions, nil
}

/**
 * get re:Invent sessions
 */
func getSessions(data url.Values) (sessions []Session, err error) {
	resp, err := http.PostForm("https://www.portal.reinvent.awsevents.com/connect/processSearchFilters.do", data)
	if err != nil {
		return sessions, err
	}
	defer resp.Body.Close()

	root, err := xmlpath.ParseHTML(resp.Body)
	if err != nil {
		return sessions, err
	}

	iterator := path.session.Iter(root)
	for iterator.Next() {
		node := iterator.Node()
		session := Session{}
		if value, ok := path.sReference.String(node); ok {
			session.ID = strings.Replace(value, "sessionDetail.ww?SESSION_ID=", "", 1)
		}
		if value, ok := path.sAbbreviation.String(node); ok {
			session.Abbreviation = strings.Replace(value, " - ", "", 1)
		}
		if value, ok := path.sTitle.String(node); ok {
			session.Title = value
		}
		if value, ok := path.sAbstract.String(node); ok {
			session.Abstract = strings.TrimSpace(value)
		}
		if value, ok := path.sLength.String(node); ok {
			session.Length = value
		}
		if value, ok := path.sType.String(node); ok {
			session.Type = value
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

/**
 * get re:Invent session details
 */
func getSessionDetails(session Session) (sesson Session, err error) {
	if session.ID == "" {
		return session, nil
	}
	resp, err := http.Get("https://www.portal.reinvent.awsevents.com/connect/sessionDetail.ww?SESSION_ID=" + session.ID)
	if err != nil {
		return session, err
	}
	defer resp.Body.Close()

	root, err := xmlpath.ParseHTML(resp.Body)
	if err != nil {
		return session, err
	}

	if value, ok := path.dTypeID.String(root); ok {
		i, _ := strconv.Atoi(value)
		session.TypeID = i
	}
	if value, ok := path.dType.String(root); ok {
		session.Type = strings.TrimSpace(value)
	}
	if value, ok := path.dTrackID.String(root); ok {
		i, _ := strconv.Atoi(value)
		session.TrackID = i
	}
	if value, ok := path.dTrack.String(root); ok {
		session.Track = strings.TrimSpace(value)
	}
	if value, ok := path.dLevelID.String(root); ok {
		i, _ := strconv.Atoi(value)
		session.LevelID = i
	}
	if value, ok := path.dLevel.String(root); ok {
		session.Level = strings.TrimSpace(value)
	}
	return session, nil
}
