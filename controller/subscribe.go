package controller

import (
	"PrayerService/model"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/lucsky/cuid"
)

// Subscribe godoc
// @Summary      User subscribes to a board
// @Description  This endpoint allows a user to subscribe to a board and receives all messages from the board
// @Tags         Boards
// @Accept       json
// @Produce      json
// @Param		 event	body	model.Event	false	"Event data"
// @Param		 Board	header		string	true	"Board ID"
// @Param		 Token	header		string	true	"Authentication header"
// @Success      200  				{object}  	model.Board
// @Failure		 400				{string}	string	"Bad Request"
// @Failure		 401				{string}	string	"Unauthorized"
// @Failure		 500				{string}	string	"Internal Server Error"
// @Router       /subscribe [get]
func (controller *Controller) Subscribe(w http.ResponseWriter, r *http.Request) {
	conn, err := controller.Upgrader.Upgrade(w, r, nil)
	log.Println("New connection:", conn.RemoteAddr().String())
	if err != nil {
		log.Println(err)
		return
	}
	boardId := r.Header.Get("Board")
	if len(boardId) == 0 {
		boardId = "1"
	}
	isDebug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	if isDebug {
		log.Println("Headers")
		for name, values := range r.Header {
			for _, value := range values {
				log.Println(name, value)
			}
		}
	}
	user, _ := r.Context().Value(UserKey).(model.User)
	boardIndex, board := controller.findBoard(boardId)
	client := &model.Client{
		ID:                 cuid.New(),
		DeviceId:           user.UserId + conn.RemoteAddr().Network(),
		BoardID:            boardId,
		IP:                 conn.RemoteAddr().String(),
		User:               user,
		CanReceiveMessages: true,
	}
	client.Send = func(messageType int, message []byte) {
		if err := conn.WriteMessage(messageType, message); err != nil {
			board.SetCanReceiveMessages(client, false)
			conn.Close()
			log.Println("client.send", err)
			return
		}
	}
	storedPrayers, _ := json.Marshal(controller.boards[boardIndex].GetBoardData())
	client.Send(1, storedPrayers)
	controller.AddClient(*client)

	go func() {
		defer func() {
			log.Println("Lost connection")
			log.Println(conn.RemoteAddr().String())
			conn.Close()
		}()
		for {
			messageType, message, err := conn.ReadMessage()
			log.Println("Received message")
			log.Println("Message Type:", messageType)
			log.Println(string(message))
			if err != nil {
				board.SetCanReceiveMessages(client, false)
				conn.Close()
				log.Println("client.receive", err)
				return
			}
			controller.broadcast(boardId, client, message, messageType)
		}
	}()
}
