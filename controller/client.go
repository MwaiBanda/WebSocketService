package controller

type Client struct {
	ID   string
	IP   string
	Send func(int, []byte)
}
