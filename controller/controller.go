package controller

import (
	constants "PrayerService/constants"
	"PrayerService/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Controller struct {
	Upgrader  websocket.Upgrader
	clients   []Client
	prayers   []model.Prayer
	broadcast func(int, []byte)
}

func GetInstance() *Controller {
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
				ID:          1,
				Title:       "Healing for my dog",
				Description: "I need prayer for my dog, he is sick",
				Comments: []model.Comment{
					{
						ID:       1,
						PrayerID: 1,
						Comment:  "He will be healed in Jesus name",
					},
				},
			},
		},
		clients: []Client{},
	}
	c.broadcast = func(messageType int, message []byte) {
		event := model.Event{}
		if err := json.Unmarshal(message, &event); err != nil {
			log.Println(err)
		}
		if event.Type == constants.PRAYER {
            if event.Action == constants.ADD {
                prayer := model.Prayer{}
                if err := json.Unmarshal([]byte(event.Data), &prayer); err != nil {
                    log.Println(err)
                }
                c.prayers = append(c.prayers, prayer)
            } else if event.Action == constants.DELETE {

            } else if event.Action == constants.UPDATE {
                
            }
            prayers, _ := json.Marshal(c.prayers)
            var waitGroup sync.WaitGroup
            for _, client := range c.clients {
                waitGroup.Add(1)
                go func() {
                    defer waitGroup.Done()
                    client.send(messageType, prayers)
                    log.Println("Broadcasting message")
                    log.Println("Sending message to client")
                    log.Println(client.deviceId)
                }()
            }
            waitGroup.Wait()
		} else if event.Type == constants.COMMENT {
            if event.Action == constants.ADD {
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
            } else if event.Action == constants.DELETE {

            } else if event.Action == constants.UPDATE {
                
            }
			prayers, _ := json.Marshal(c.prayers)
            var waitGroup sync.WaitGroup
			for _, client := range c.clients {
                waitGroup.Add(1)
                go func() {
                    client.send(messageType, prayers)
                    log.Println("Broadcasting message")
                    log.Println("Sending message to client")
                    log.Println(client.deviceId)
                }()
			}
            waitGroup.Wait()
		}
	}
	return c
}

func (controller *Controller) AddClient(client Client) {
	filtered := []Client{}
	for _, c := range controller.clients {
		if c.deviceId != client.deviceId {
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
