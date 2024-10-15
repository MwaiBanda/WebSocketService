package model

type Prayer struct {
	ID          int `json:"id"`
	BoardID          int `json:"boardId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Comments    []Comment `json:"comments"`
	User		User `json:"user"`
}

func (u *Prayer) SetUser(user User) {
	u.User = user
  }