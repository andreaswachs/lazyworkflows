package response

import (
	"encoding/json"
)

type Workflow struct {
	Id        json.Number
	NodeId    string `json:"node_id"`
	Name      string
	Path      string
	State     string
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Url       string
	HtmlUrl   string `json:"html_url"`
	BadgeUrl  string `json:"badge_url"`
}

type Get struct {
	Workflow
}

type List struct {
	TotalCount int `json:"total_count"`
	Workflows  []Workflow
}

type Enable struct {
	Status int
}

type Disable struct {
	Status int
}

type Dispatch struct {
	Status int
}

func FromString[T any](response string, out T) error {
	return json.Unmarshal([]byte(response), &out)
}
