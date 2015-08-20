package models

type sessionType struct {
	SearchID int
	ResultID int
	Caption  string
}

var sessionTypes []sessionType

func init() {
	sessionTypes = []sessionType{
		sessionType{2, 1040, "Breakout Session"},
		sessionType{1140, 12223, "General Activity"},
		sessionType{1160, 12523, "Meal"},
		sessionType{1183, 12603, "Hands-on Lab"},
		sessionType{1184, 12604, "Certification"},
	}
}
