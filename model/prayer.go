package model

type Prayer struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PrayerTime  string `json:"prayer_time"`
	Comments    []Comment `json:"comments"`
}

