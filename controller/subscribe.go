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

func (controller *Controller) Subscribe(w http.ResponseWriter, r *http.Request) {
	conn, err := controller.Upgrader.Upgrade(w, r, nil)
	log.Println("New connection:", conn.RemoteAddr().String())
	if err != nil {
		log.Println(err)
		return
	}
	deviceId := r.Header.Get("Device_id")
	boardId := r.Header.Get("Board")
	if len(deviceId) == 0 {
		deviceId = cuid.New()
	}
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
	log.Println("User:", user)
	client := &model.Client{
		ID: deviceId,
		BoardID: boardId,
		IP: conn.RemoteAddr().String(),
	}
	client.Send = func(messageType int, message []byte) {
		if err := conn.WriteMessage(messageType, message); err != nil {
			controller.RemoveClient(*client)
			conn.Close()
			log.Println("client.send", err)
			return
		}
	}
	boardIndex, _ := controller.findBoard(boardId)
	storedPrayers, _ := json.Marshal(controller.boards[boardIndex].Prayers)
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
				controller.RemoveClient(*client)
				conn.Close()
				log.Println(err)
				return
			}
			controller.broadcast(boardId, user, message, messageType)
		}
	}()
}
