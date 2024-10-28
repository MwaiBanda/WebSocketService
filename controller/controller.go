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
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/machinebox/graphql"
)

type Controller struct {
	Upgrader   websocket.Upgrader
	boards    []model.Board
	broadcast  func(string, *model.Client, []byte, int)
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
				}, Clients: []model.Client{}, Boards: []model.Board{
					{
						ID:    1,
						Title: "General",
					},
					{
						ID:    2,
						Title: "Healing",
					},
				},
			},
			{
				ID:    2,
				Title: "Healing",
				Prayers: []model.Prayer{
					{
						ID:          100,
						BoardID:          2,
						Title:       "Healing for my cat",
						Description: "I need prayer for my cat, he is sick",
						User: model.User{
							UserId:   "1",
							FirstName: "Mary",
							LastName:  "Jane",
							UserName: "maryjane",
							ScreenName: "Mary Jane",
						},
						Comments: []model.Comment{
							{
								ID:       100,
								PrayerID: 100,
								Comment:  "He will be healed in Jesus name",
								User: model.User{
									UserId:   "1",
									FirstName: "Mary",
									LastName:  "Jane",
									UserName: "maryjane",
									ScreenName: "Mary Jane",
								},
							},
						},
					},
				}, Clients: []model.Client{}, Boards: []model.Board{
					{
						ID:    1,
						Title: "General",
					},
					{
						ID:    2,
						Title: "Healing",
					},
				},
			},
		},
	}
	c.broadcast = func(
		boardId string,
		client *model.Client,
		message []byte,
		messageType int,
	) {
		event := model.Event{}
		boardIndex, currentBoard := c.findBoard(boardId)
		if err := json.Unmarshal(message, &event); err != nil {
			log.Println(err)
		}
		if event.Type == constants.PRAYER {
			if event.Action == constants.ADD {
				prayer := model.Prayer{}
				if err := json.Unmarshal([]byte(event.Data), &prayer); err != nil {
					log.Println(err)
				}
				prayer.SetUser(client.User)
				c.boards[boardIndex].Prayers = append(c.boards[boardIndex].Prayers, prayer)
			} else if event.Action == constants.DELETE {
				prayer := model.Prayer{}
				if err := json.Unmarshal([]byte(event.Data), &prayer); err != nil {
					log.Println(err)
				}
				for i, p := range c.boards[boardIndex].Prayers {
					if p.ID == prayer.ID {
						c.boards[boardIndex].Prayers = append(c.boards[boardIndex].Prayers[:i], c.boards[boardIndex].Prayers[i+1:]...)
						break
					}
				}
			} else if event.Action == constants.UPDATE {

			}
			board, _ := json.Marshal(c.boards[boardIndex].GetBoardData())
			for _, client := range c.boards[boardIndex].Clients {
				if client.BoardID == boardId {
					client.Send(messageType, board)
					log.Println("Broadcasting message")
					log.Println("Sending message to client")
					log.Println(client.User.UserName)
				}
			}
			log.Println("Number of clients:", len(c.boards[boardIndex].Clients))
		} else if event.Type == constants.COMMENT {
			if event.Action == constants.ADD {
				comment := model.Comment{}
				if err := json.Unmarshal([]byte(event.Data), &comment); err != nil {
					log.Println(err)
				}
				comment.SetUser(client.User)
				for i, p := range c.boards[boardIndex].Prayers {
					if p.ID == comment.PrayerID {
						c.boards[boardIndex].Prayers[i].Comments = append(c.boards[boardIndex].Prayers[i].Comments, comment)
						break
					}
				}
			} else if event.Action == constants.DELETE {

			} else if event.Action == constants.UPDATE {

			}
			board, _ := json.Marshal(c.boards[boardIndex].GetBoardData())
			for _, client := range c.boards[boardIndex].Clients {
				if client.CanReceiveMessages {
					client.Send(messageType, board)
					log.Println("Broadcasting message")
					log.Println("Sending message to client")
					log.Println(client.User.UserName)
				}
			}
		} else if event.Type == constants.BOARD {
			if event.Action == constants.UPDATE {
				board := model.BoardEvent{}
				if err := json.Unmarshal([]byte(event.Data), &board); err != nil {
					log.Println(err)
				}

				log.Println("Switching to Board ID:", board.ID)
				newBoardId := strconv.Itoa(board.ID)
				_, newBoard := c.findBoard(newBoardId)
				c.MoveClient(newBoardId, client)
				log.Println("New board:", newBoard.Clients, "Old board:", currentBoard.Clients)
				boardJson, _ := json.Marshal(newBoard.GetBoardData())
				client.Send(messageType, boardJson)
			}
		}
	}
	return c
}

func (controller *Controller) findBoard(boardId string) (int, model.Board) {
	id, _ := strconv.Atoi(boardId)
	for i, board := range controller.boards {
		if board.ID == id {
			return i, board
		}
	}
	return 0, model.Board{}
}
func (controller *Controller) AddClient(client model.Client) {
	filtered := []model.Client{}
	boardIndex, board := controller.findBoard(client.BoardID)
	for _, c := range controller.boards[boardIndex].Clients {
		if c.ID != client.ID {
			filtered = append(filtered, c)

		}
	}
	controller.boards[boardIndex].Clients = filtered
	client.SetBoardId(strconv.Itoa(board.ID))
	controller.boards[boardIndex].Clients = append(controller.boards[boardIndex].Clients, client)
	log.Println("Number of clients:", len(controller.boards[boardIndex].Clients))
}

func (controller *Controller) MoveClient(newBoardId string, client *model.Client) {
	newBoardIndex, _ := controller.findBoard(newBoardId)
	boardIndex, board := controller.findBoard(client.BoardID)
	for _, c := range controller.boards[boardIndex].Clients {
		if c.ID == client.ID {
			newClient := client
			newClient.SetBoardId(newBoardId)
			client.SetBoardId(newBoardId)
			newClient.SetCanReceiveMessages(true)
			controller.boards[newBoardIndex].AddClient(*newClient)
			controller.boards[newBoardIndex].SetCanReceiveMessages(newClient, true)
			board.SetCanReceiveMessages(client, false)
			break
		}
	}
}

func (controller *Controller) RemoveClient(client model.Client) {
	boardIndex, _ := controller.findBoard(client.BoardID)
	for _, c := range controller.boards[boardIndex].Clients {
		if c.ID == client.ID {
			controller.boards[boardIndex].RemoveClient(client)
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
