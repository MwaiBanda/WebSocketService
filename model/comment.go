package model

type Comment struct {
	ID        string `json:"id"`
	PrayerID  string `json:"prayer_id"`
	Comment   string `json:"comment"`
}