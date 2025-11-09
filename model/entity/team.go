package entity

type Team struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	IsPublic    bool     `json:"ispublic"`
	UsersIds    []string `json:"users"`
}

func NewTeam(id, name, desc string, isPublic bool, Users []string) *Team {
	return &Team{
		Id:          id,
		Name:        name,
		Description: desc,
		IsPublic:    isPublic,
		UsersIds: func() []string {
			if Users == nil {
				return []string{}
			}
			return Users
		}(),
	}
}
