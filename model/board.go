package model

type Board struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Prayers    []Prayer  `json:"prayers"`
	Boards	 []Board   `json:"boards"`
}