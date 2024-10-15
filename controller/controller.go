package controller

import (
	constants "PrayerService/constants"
	"PrayerService/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/machinebox/graphql"
)

type Controller struct {
	Upgrader   websocket.Upgrader
	clients    []Client
	prayers    []model.Prayer
	broadcast  func(model.User, []byte, int)
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
				User: model.User{
					UserId:   "1",
					FirstName: "John",
					LastName:  "Doe",
					UserName: "johndoe",
					ScreenName: "John Doe",
				},
				Comments: []model.Comment{
					{
						ID:       1,
						PrayerID: 1,
						Comment:  "He will be healed in Jesus name",
						User: model.User{
							UserId:   "1",
							FirstName: "John",
							LastName:  "Doe",
							UserName: "johndoe",
							ScreenName: "John Doe",
						},
					},
				},
			},
		},
		clients: []Client{},
	}
	c.broadcast = func(user model.User, message []byte, messageType int) {
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
				prayer.SetUser(user)
				c.prayers = append(c.prayers, prayer)
			} else if event.Action == constants.DELETE {
				prayer := model.Prayer{}
				if err := json.Unmarshal([]byte(event.Data), &prayer); err != nil {
					log.Println(err)
				}
				for i, p := range c.prayers {
					if p.ID == prayer.ID {
						c.prayers = append(c.prayers[:i], c.prayers[i+1:]...)
						break
					}
				}
			} else if event.Action == constants.UPDATE {

			}
			prayers, _ := json.Marshal(c.prayers)
			for _, client := range c.clients {
				client.Send(messageType, prayers)
				log.Println("Broadcasting message")
				log.Println("Sending message to client")
				log.Println(user.UserName)
			}
		} else if event.Type == constants.COMMENT {
			if event.Action == constants.ADD {
				comment := model.Comment{}
				if err := json.Unmarshal([]byte(event.Data), &comment); err != nil {
					log.Println(err)
				}
				comment.SetUser(user)
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
			for _, client := range c.clients {
				client.Send(messageType, prayers)
				log.Println("Broadcasting message")
				log.Println("Sending message to client")
				log.Println(user.UserName)
			}
		}
	}
	return c
}

func (controller *Controller) AddClient(client Client) {
	filtered := []Client{}
	for _, c := range controller.clients {
		if c.ID != client.ID {
			filtered = append(filtered, c)
		}
	}
	controller.clients = filtered
	controller.clients = append(controller.clients, client)
	log.Println("Number of clients:", len(controller.clients))
}

func (controller *Controller) RemoveClient(clientId string) {
	for _, c := range controller.clients {
		if c.ID != clientId {
			controller.clients = append(controller.clients, c)
			break
		}
	}
}

type contextKey string

const UserKey contextKey = "user"

func (controller *Controller) Auth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Token")
		if len(token) == 0 {
			log.Println("Unauthorized", r.URL.Path)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Println("Token:", token)
		log.Println(r.URL.Path)
		user, err := controller.getProfile(token)
		if err != nil {
			log.Println("Unauthorized", r.URL.Path)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		req := r.WithContext(context.WithValue(r.Context(), UserKey, user))
		f(w, req)
	}
}

func (controller *Controller) getProfile(token string) (model.User, error) {
	client := graphql.NewClient(os.Getenv("GRAPHQL_URL"))
	req := graphql.NewRequest(`
		query Query {
            getUserProfile {
                lastName
                firstName
                email
                screenName
                userName
                userId
            }
        }
	`)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	ctx := context.Background()
	var respData model.UserResponse
	if err := client.Run(ctx, req, &respData); err != nil {
		return model.User{}, fmt.Errorf("failed to get user profile: %v", err)
	}
	return respData.GetUserProfile, nil
}
