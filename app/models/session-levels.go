package models

// SessionLevel represents re:Invent's session track
type SessionLevel struct {
	ID      int
	Caption string
}

var sessionLevels []SessionLevel

func init() {
	sessionLevels = []SessionLevel{
		SessionLevel{10141, "Introductory (200 level)"},
		SessionLevel{10142, "Advanced (300 level)"},
		SessionLevel{10143, "Expert (400 level)"},
	}
}

// SessionLevels returns all session types
func SessionLevels() []SessionLevel {
	return sessionLevels
}
