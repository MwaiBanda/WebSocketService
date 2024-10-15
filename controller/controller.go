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
	boards    []model.Board
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
		boards: []model.Board{
			{
				ID:    1,
				Title: "General",
				Prayers: []model.Prayer{
					{
						ID:          1,
						BoardID:          1,
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
				}, Clients: []model.Client{}, Boards: []model.Board{},
			},
		},
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
				c.boards[0].Prayers = append(c.boards[0].Prayers, prayer)
			} else if event.Action == constants.DELETE {
				prayer := model.Prayer{}
				if err := json.Unmarshal([]byte(event.Data), &prayer); err != nil {
					log.Println(err)
				}
				for i, p := range c.boards[0].Prayers {
					if p.ID == prayer.ID {
						c.boards[0].Prayers = append(c.boards[0].Prayers[:i], c.boards[0].Prayers[i+1:]...)
						break
					}
				}
			} else if event.Action == constants.UPDATE {

			}
			prayers, _ := json.Marshal(c.boards[0].Prayers)
			for _, client := range c.boards[0].Clients {
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
				for i, p := range c.boards[0].Prayers {
					if p.ID == comment.PrayerID {
						c.boards[0].Prayers[i].Comments = append(c.boards[0].Prayers[i].Comments, comment)
						break
					}
				}
			} else if event.Action == constants.DELETE {

			} else if event.Action == constants.UPDATE {

			}
			prayers, _ := json.Marshal(c.boards[0].Prayers)
			for _, client := range c.boards[0].Clients {
				client.Send(messageType, prayers)
				log.Println("Broadcasting message")
				log.Println("Sending message to client")
				log.Println(user.UserName)
			}
		}
	}
	return c
}

func (controller *Controller) AddClient(client model.Client) {
	filtered := []model.Client{}
	for _, c := range controller.boards[0].Clients {
		if c.ID != client.ID {
			filtered = append(filtered, c)
		}
	}
	controller.boards[0].Clients = filtered
	controller.boards[0].Clients = append(controller.boards[0].Clients, client)
	log.Println("Number of clients:", len(controller.boards[0].Clients))
}

func (controller *Controller) RemoveClient(clientId string) {
	for _, c := range controller.boards[0].Clients {
		if c.ID != clientId {
			controller.boards[0].Clients = append(controller.boards[0].Clients, c)
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
