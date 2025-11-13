package entity

type VoiceRoom struct {
	Id           string   `json:"id"`
	TeamId       string   `json:"teamId"`
	Name         string   `json:"name"`
	IsActive     bool     `json:"isActive"`
	Participants []string `json:"participants"`
	CreatedBy    string   `json:"createdBy"`
	CreatedAt    int64    `json:"createdAt"`
}
