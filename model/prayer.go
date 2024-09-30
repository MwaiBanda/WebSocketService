package model

type Prayer struct {
	ID          int `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Comments    []Comment `json:"comments"`
}

