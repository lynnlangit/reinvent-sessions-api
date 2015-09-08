package models

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/supinf/reinvent-sessions-api/app/aws"
	"github.com/supinf/reinvent-sessions-api/app/config"
	"github.com/supinf/reinvent-sessions-api/app/logs"
)

var sessionTable string

// Session represents re:Invent's session
type Session struct {
	ID           string `json:"id"`
	Abbreviation string `json:"abbreviation"`
	Title        string `json:"title"`
	Abstract     string `json:"abstract"`
	// Recommended  bool   `json:"recommended"`
	// Speakers     []string  `json:"speakers,omitempty"`
	Date    int64     `json:"date"`
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Room    string    `json:"room"`
	Length  string    `json:"length"`
	TypeID  int       `json:"typeId"`
	Type    string    `json:"type"`
	TrackID int       `json:"trackId"`
	Track   string    `json:"track"`
	LevelID int       `json:"levelId"`
	Level   string    `json:"level"`
}

// Sessions is a type of Session slice
type Sessions []Session

var sessionOnce sync.Once

func init() {
	sessionOnce.Do(func() {
		r, _ := regexp.Compile("[^a-zA-Z0-9_\\.]")
		sessionTable = r.ReplaceAllString(strings.ToLower(config.NewConfig().Name), "-") + "-sessions"
		if found, _ := comfirmSessionTableExists(); !found {
			go func() {
				time.Sleep(10 * time.Second)
				SyncReInventSessions(true)
			}()
		}
	})
}

// GetSessions lists all sessions from DynamoDB
//  @return sessions []models.Session
func GetSessions() (sessions Sessions, count int64, err error) {
	records, count, err := aws.DynamoRecords(sessionTable)
	if err != nil {
		return sessions, 0, err
	}
	return toSessions(records), count, nil
}

// GetSession retrives a specified session from DynamoDB
//  @param  id string
//  @return session models.Session
func GetSession(id string) (session Session, err error) {
	record, err := aws.DynamoRecord(sessionTable, map[string]*dynamodb.AttributeValue{
		"ID": aws.DynamoAttributeS(id),
	})
	if err != nil {
		return session, nil
	}
	return toSession(record), nil
}

// cast DynamoDB records to Sessions
func toSessions(records []map[string]*dynamodb.AttributeValue) (sessions Sessions) {
	for _, record := range records {
		session := toSession(record)
		sessions = append(sessions, session)
	}
	if len(sessions) == 0 {
		sessions = make([]Session, 0)
	}
	sort.Sort(sessions)
	return sessions
}

// cast DynamoDB record to a Session
func toSession(record map[string]*dynamodb.AttributeValue) Session {
	session := Session{}
	session.ID = aws.DynamoS(record, "ID")
	session.Abbreviation = aws.DynamoS(record, "Abbreviation")
	session.Title = aws.DynamoS(record, "Title")
	session.Abstract = aws.DynamoS(record, "Abstract")
	// session.Recommended = aws.DynamoB(record, "Recommended")

	// FIXME http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_AttributeValue.html
	// session.Speakers = []string{}
	session.Date = aws.DynamoN64(record, "Date")
	session.Start = aws.DynamoD(record, "Start")
	session.End = aws.DynamoD(record, "End")
	session.Room = aws.DynamoS(record, "Room")
	session.Length = aws.DynamoS(record, "Length")
	session.TypeID = aws.DynamoN(record, "TypeId")
	session.Type = aws.DynamoS(record, "Type")
	session.TrackID = aws.DynamoN(record, "TrackId")
	session.Track = aws.DynamoS(record, "Track")
	session.LevelID = aws.DynamoN(record, "LevelId")
	session.Level = aws.DynamoS(record, "Level")
	return session
}

// Contains checks whether it contains key word or not
func (s *Session) Contains(words []string) bool {
	session := s.toUpperFieldSession()
	match := true
	for _, word := range words {
		match = match && (strings.Contains(session.ID, word) ||
			strings.Contains(session.Abbreviation, word) ||
			strings.Contains(session.Title, word) ||
			strings.Contains(session.Abstract, word) ||
			strings.Contains(session.Room, word) ||
			strings.Contains(session.Type, word) ||
			strings.Contains(session.Track, word) ||
			strings.Contains(session.Level, word))
	}
	return match
}

func (s *Session) toUpperFieldSession() Session {
	session := Session{}
	session.ID = s.ID
	session.Abbreviation = strings.ToUpper(s.Abbreviation)
	session.Title = strings.ToUpper(s.Title)
	session.Abstract = strings.ToUpper(s.Abstract)
	session.Room = strings.ToUpper(s.Room)
	session.Type = strings.ToUpper(s.Type)
	session.Track = strings.ToUpper(s.Track)
	session.Level = strings.ToUpper(s.Level)
	return session
}

func (s Session) persist() error {
	items := map[string]*dynamodb.AttributeValue{}
	items["ID"] = aws.DynamoAttributeS(s.ID)
	if s.Abbreviation != "" {
		items["Abbreviation"] = aws.DynamoAttributeS(s.Abbreviation)
	}
	if s.Title != "" {
		items["Title"] = aws.DynamoAttributeS(s.Title)
	}
	if s.Abstract != "" {
		items["Abstract"] = aws.DynamoAttributeS(s.Abstract)
	}
	// items["Recommended"] = aws.DynamoAttributeB(s.Recommended)
	// session.Speakers = []string{}
	items["Date"] = aws.DynamoAttributeN64(s.Date)
	items["Start"] = aws.DynamoAttributeD(s.Start)
	items["End"] = aws.DynamoAttributeD(s.End)
	if s.Room != "" {
		items["Room"] = aws.DynamoAttributeS(s.Room)
	}
	if s.Length != "" {
		items["Length"] = aws.DynamoAttributeS(s.Length)
	}
	items["TypeId"] = aws.DynamoAttributeN(s.TypeID)
	if s.Type != "" {
		items["Type"] = aws.DynamoAttributeS(s.Type)
	}
	items["TrackId"] = aws.DynamoAttributeN(s.TrackID)
	if s.Track != "" {
		items["Track"] = aws.DynamoAttributeS(s.Track)
	}
	items["LevelId"] = aws.DynamoAttributeN(s.LevelID)
	if s.Level != "" {
		items["Level"] = aws.DynamoAttributeS(s.Level)
	}
	_, err := aws.DynamoPutItem(sessionTable, items)
	return err
}

func comfirmSessionTableExists() (found bool, err error) {
	if _, err := aws.DynamoTable(sessionTable); err == nil {
		return true, nil
	}
	logs.Debug.Print("[model] Session table was not found. Try to make it. @aws.DynamoCreateTable")
	attributes := map[string]string{
		"ID": "S",
	}
	keys := map[string]string{
		"ID": "HASH",
	}
	_, err = aws.DynamoCreateTable(sessionTable, attributes, keys, 1, 1)
	return false, err
}

func (s Sessions) Len() int {
	return len(s)
}

func (s Sessions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Sessions) Less(i, j int) bool {
	a, _ := strconv.Atoi(s[i].ID)
	b, _ := strconv.Atoi(s[j].ID)
	return a < b
}
