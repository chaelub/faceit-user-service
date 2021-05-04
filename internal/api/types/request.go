package types

type (
	UserRequest struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Nickname  string `json:"nickname"`
		Password  string `json:"password,omitempty"`
		Email     string `json:"email"`
		Country   string `json:"country"`
	}
)
