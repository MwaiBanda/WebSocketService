package controller

import "github.com/gorilla/websocket"

type Client struct {
	
	controller *Controller

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send func(int, []byte)
}
