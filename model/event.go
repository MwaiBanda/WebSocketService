package model

type Event struct {
	Type string `json:"type"`
	Action string `json:"action"`
	Data string `json:"data"`
}