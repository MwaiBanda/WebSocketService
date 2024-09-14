package controller

import (
	"PrayerService/model"
	"net/http"

	"github.com/gorilla/websocket"
)

type Controller struct {
	Upgrader websocket.Upgrader
	clients  []Client
    prayers []model.Prayer
	broadcast func(int, []byte)
}

func NewController() *Controller {
	c := &Controller{
        Upgrader: websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
            CheckOrigin: func(r *http.Request) bool {
                return true
            },
        },
        clients: []Client{},
    }
    c.broadcast = func(messageType int, message []byte) {
        for _, client := range c.clients {
            client.send(messageType, message)
        }
    }
    return c
}

func (controller *Controller) AddClient(client Client) {
	controller.clients = append(controller.clients, client)
}

func (controller *Controller) RemoveClient(client Client) {
    for i, c := range controller.clients {
        if c.id == client.id {
            controller.clients = append(controller.clients[:i], controller.clients[i+1:]...)
            break
        }
    }
}

