package controller

import "github.com/gorilla/websocket"

type Client struct {
	id string
	controller *Controller
	conn *websocket.Conn
	// Buffered func of outbound messages.
	send func(int, []byte)
}
