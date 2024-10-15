package model

type Client struct {
	ID   string
	BoardID   string
	IP   string
	Send func(int, []byte)
}
