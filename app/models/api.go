// Package models defines data models and related functions
package models

// API represents api intercace
type API struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Method      string         `json:"method"`
	Parameters  []APIParameter `json:"parameters"`
	Example     string         `json:"example"`
}

// APIParameter represents api parameter
type APIParameter struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Necessary   bool   `json:"necessary"`
}

// ListAPI returns all API
func ListAPI() (apis []API) {
	apis = append(apis, API{
		Name:        "/reinvent/sessions",
		Description: "Lists re:Invent sessions",
		Method:      "GET",
		Parameters: []APIParameter{
			APIParameter{
				Key:         "output",
				Description: "The formatting style: html | json",
				Necessary:   false,
			},
			APIParameter{
				Key:         "id",
				Description: "Session ID",
				Necessary:   false,
			},
			APIParameter{
				Key:         "q",
				Description: "Space seperated words to filter the response data",
				Necessary:   false,
			},
		},
		Example: "reinvent/sessions?output=json&q=security%20400",
	})
	return apis
}
