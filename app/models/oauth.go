package models

// OAuthAuthorizedToken represents authorized person
type OAuthAuthorizedToken struct {
	ID                string `json:"-"`
	ScreenName        string `json:"name"`
	AccessTokenKey    string `json:"token"`
	AccessTokenSecret string `json:"secret"`
	CognitoPoolID     string `json:"pool"`
	CognitoRoleArn    string `json:"role"`
}
