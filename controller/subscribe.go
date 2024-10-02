package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

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
	if len(deviceId) == 0 {
		deviceId = cuid.New()
	}

	if os.Getenv("DEBUG") == "true" {
		log.Println("Headers")
		for name, values := range r.Header {
			for _, value := range values {
				log.Println(name, value)
			}
		}
	}
	
	client := &Client{
		ID: deviceId,
		IP: conn.RemoteAddr().String(),
		Send: func(messageType int, message []byte) {
			if err := conn.WriteMessage(messageType, message); err != nil {
				controller.RemoveClient(deviceId)
				conn.Close()
				log.Println("client.send", err)
				return
			}
		},
	}
	storedPrayers, _ := json.Marshal(controller.prayers)
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
				controller.RemoveClient(deviceId)
				conn.Close()
				log.Println(err)
				return
			}
			controller.broadcast(messageType, message)
		}
	}()
}
