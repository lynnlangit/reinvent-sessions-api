package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/supinf/reinvent-sessions-api/app/misc"
)

type sessionSchedule struct {
	Action        string    `json:"-"`
	ImageStyle    string    `json:"-"`
	Message       string    `json:"message"`
	Date          int64     `json:"_"`
	StartTime     string    `json:"startTime"`
	StartDateTime time.Time `json:"_"`
	EndTime       string    `json:"endTime"`
	EndDateTime   time.Time `json:"_"`
	Room          string    `json:"room"`
	RoomID        int       `json:"-"`
	MapID         int       `json:"-"`
	SessionTimeID int       `json:"-"`
}

type SessionScheduleList struct {
	Data []sessionSchedule
}

func SessionSchedule(id string, cookies []*http.Cookie, sessionID string) (schedule SessionScheduleList, err error) {
	cookies = append(cookies, &http.Cookie{
		Name:  "DWRSESSIONID",
		Value: sessionID,
		Path:  "/",
	})
	url := "https://www.portal.reinvent.awsevents.com/connect/dwr/call/plaincall/ConnectAjax.getSchedulingJSON.dwr"
	values := dwrQueryStrings("ConnectAjax", "getSchedulingJSON", "number:"+id, "/connect/sessionDetail.ww", sessionID+"/"+dwrPageID())
	client, req, err := dwrHTTPClient("POST", url, values, cookies)
	if err != nil {
		return schedule, err
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	response := strings.Replace(strings.Replace(string(body), "\n", "", -1), "\\", "", -1)
	trimed := regexp.MustCompile(`^.*handleCallback\("0","0","`).ReplaceAllString(response, "")
	trimed = regexp.MustCompile(`"\)\;\}\)\(\)\;$`).ReplaceAllString(trimed, "")

	temp := SessionScheduleList{}
	err = json.Unmarshal([]byte(trimed), &temp)
	if err != nil {
		return schedule, err
	}
	now := time.Now()
	for _, sche := range temp.Data {
		s, _ := time.Parse("Monday, Jan 2, 3:4 PM", sche.StartTime)
		e, _ := time.Parse("3:4 PM", sche.EndTime)
		schedule.Data = append(schedule.Data, sessionSchedule{
			Action:        sche.Action,
			ImageStyle:    sche.ImageStyle,
			Message:       sche.Message,
			Date:          time.Date(now.Year(), s.Month(), s.Day(), 0, 0, 0, 0, time.UTC).Unix(),
			StartTime:     sche.StartTime,
			StartDateTime: time.Date(now.Year(), s.Month(), s.Day(), s.Hour(), s.Minute(), 0, 0, time.UTC),
			EndTime:       sche.EndTime,
			EndDateTime:   time.Date(now.Year(), s.Month(), s.Day(), e.Hour(), e.Minute(), 0, 0, time.UTC),
			Room:          sche.Room,
			RoomID:        sche.RoomID,
			MapID:         sche.MapID,
			SessionTimeID: sche.SessionTimeID,
		})
	}
	return schedule, nil
}

func DwrSession(cookies *[]*http.Cookie) error {
	data := url.Values{}
	data["key"] = []string{"Meal"}
	data["searchType"] = []string{"session"}
	data["sessionTypeID"] = []string{fmt.Sprint(1160)}
	resp, err := http.PostForm("https://www.portal.reinvent.awsevents.com/connect/processSearchFilters.do", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	for _, cookie := range resp.Cookies() {
		*cookies = append(*cookies, cookie)
	}
	return nil
}

func DwrSessionID(cookies []*http.Cookie) string {
	url := "https://www.portal.reinvent.awsevents.com/connect/dwr/call/plaincall/__System.generateId.dwr"
	values := dwrQueryStrings("__System", "generateId", "", "/connect/search.ww", "")
	client, req, err := dwrHTTPClient("POST", url, values, cookies)
	if err != nil {
		return ""
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	response := strings.Replace(strings.Replace(string(body), "\n", "", -1), "\"", "", -1)

	reg := regexp.MustCompile(`,0,(\S+)\)\;\}\)\(\)\;$`)
	matches := reg.FindStringSubmatch(response)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func dwrQueryStrings(script, method, param, page, sessionID string) *url.Values {
	values := url.Values{}
	values.Add("callCount", "1")
	values.Add("c0-scriptName", script)
	values.Add("c0-methodName", method)
	values.Add("c0-id", "0")
	if !misc.ZeroOrNil(param) {
		values.Add("c0-param0", param)
	}
	values.Add("batchId", "0")
	values.Add("instanceId", "0")
	values.Add("page", page)
	values.Add("scriptSessionId", sessionID)
	values.Add("windowName", "")
	return &values
}

func dwrHTTPClient(method, uri string, values *url.Values, cookies []*http.Cookie) (client *http.Client, req *http.Request, err error) {
	req, err = http.NewRequest(method, uri, strings.NewReader(values.Encode()))
	if err != nil {
		return
	}
	cookieURL, _ := url.Parse(uri)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(cookieURL, cookies)

	return &http.Client{Jar: jar}, req, nil
}

func dwrTokenify(number int64) string {
	charmap := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ*$"
	remainder := float64(number)
	result := ""
	for remainder > 0 {
		result += string(charmap[int64(remainder)&0x3F])
		remainder = math.Floor(remainder / 64)
	}
	return result
}

func dwrPageID() string {
	return dwrTokenify(time.Now().Unix()) + "-" + dwrTokenify(int64(rand.Intn(1E16)))
}
