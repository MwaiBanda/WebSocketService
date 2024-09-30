package controller

import (
	"PrayerService/model"
	"encoding/json"
	"fmt"
	"log"
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
        prayers: []model.Prayer{
            {
                ID: 1,
                Title: "Healing for my dog",
                Description: "I need prayer for my dog, he is sick",
                Comments: []model.Comment{
                    {
                        ID: 1,
                        PrayerID: 1,
                        Comment: "He will be healed in Jesus name",
                    },
                },
            },
        },
        clients: []Client{},
    }
    c.broadcast = func(messageType int, message []byte) {
        log.Println("Broadcasting message")
        for _, client := range c.clients {
            log.Println("Sending message to client")
            log.Println(client.id)
            log.Println(client.deviceId)
           
            event := model.Event{}
            if err := json.Unmarshal(message, &event); err != nil {
                log.Println(err)
            }
            if event.Type == "prayer" {
                prayer := model.Prayer{}
                if err := json.Unmarshal([]byte(event.Data), &prayer); err != nil {
                    log.Println(err)
                }
                c.prayers = append(c.prayers, prayer)
                prayers, _ := json.Marshal(c.prayers)
                client.send(messageType, prayers)
            } else if event.Type == "comment" {
                comment := model.Comment{}
                if err := json.Unmarshal([]byte(event.Data), &comment); err != nil {
                    log.Println(err)
                }
                for i, p := range c.prayers {
                    if p.ID == comment.PrayerID {
                        c.prayers[i].Comments = append(c.prayers[i].Comments, comment)
                        break
                    }
                }
                prayers, _ := json.Marshal(c.prayers)
                client.send(messageType, prayers)
            }
        }
    }
    return c
}

func (controller *Controller) AddClient(client Client) {
           
    filtered := []Client{}
    for _, c := range controller.clients {
        if(c.deviceId != client.deviceId) {
            filtered = append(filtered, c)
        }
    }
    controller.clients = filtered
	controller.clients = append(controller.clients, client)
    fmt.Println("Number of clients:", len(controller.clients))
}

func (controller *Controller) RemoveClient(clientId string) {
    for _, c := range controller.clients {
        if c.id != clientId {
             controller.clients = append(controller.clients, c)
            break
        }
    }
}

