package model

type Comment struct {
	ID        int `json:"id"`
	PrayerID  int `json:"prayer_id"`
	Comment   string `json:"comment"`
}