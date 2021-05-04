package models

import "encoding/json"

type (
	publicUser User

	User struct {
		Id        int64  `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Nickname  string `json:"nickname"`
		Password  string `json:"password,omitempty"`
		Email     string `json:"email"`
		Country   string `json:"country"`
	}
)

func (u User) MarshalJSON() ([]byte, error) {
	u.Password = ""
	return json.Marshal(publicUser(u))
}
