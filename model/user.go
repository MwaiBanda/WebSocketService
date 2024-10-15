package model

type UserResponse struct {
	GetUserProfile User `json:"getUserProfile"`
}

type User struct {
	Email    string `json:"email"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	ScreenName string `json:"screenName"`
	UserName string `json:"userName"`
	UserId string `json:"userId"`
}