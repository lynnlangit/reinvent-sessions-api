package models

type sessionTrack struct {
	ID      int
	Caption string
}

var sessionTracks []sessionTrack

func init() {
	sessionTracks = []sessionTrack{
		sessionTrack{10481, "DevOps"},
		sessionTrack{10482, "Architecture"},
		sessionTrack{10483, "Big Data & Analytics"},
		sessionTrack{10484, "Compute"},
		sessionTrack{10485, "Databases"},
		sessionTrack{10486, "Developer Tools"},
		sessionTrack{10487, "Storage & Content Delivery"},
		sessionTrack{10488, "Spotlight"},
		sessionTrack{10489, "Security & Compliance"},
		sessionTrack{10490, "Networking"},
		sessionTrack{10491, "IT Strategy & Migration"},
		sessionTrack{10492, "Gaming"},
		sessionTrack{10493, "Mobile Developer & IoT"},
	}
}
