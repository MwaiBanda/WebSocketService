package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lucsky/cuid"
)

func (controller *Controller) Subscribe(w http.ResponseWriter, r *http.Request) {
	conn, err := controller.Upgrader.Upgrade(w, r, nil)
	fmt.Println("New connection")
	fmt.Println(conn.RemoteAddr().String())
	if err != nil {
		log.Println(err)
		return
	}
	clientId := cuid.New()
	deviceId := r.Header.Get("Device_id")
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Println(name, value)
		}
	}
	client := &Client{
		id: clientId,
		deviceId: deviceId,
		controller: controller,
		conn: conn,
		send: func(messageType int, message []byte) {
			if err := conn.WriteMessage(messageType, message); err != nil {
				// controller.RemoveClient(clientId)
				conn.Close()
				log.Println("client.send", err)
				return
			}
		},
	}
	storedPrayers, _ := json.Marshal(controller.prayers)
	client.send(1, storedPrayers)
	controller.AddClient(*client)

	go func() {
		defer func() {
			fmt.Println("Lost connection")
			fmt.Println(conn.RemoteAddr().String())
			conn.Close()
		}()
		for {
			messageType, message, err := conn.ReadMessage()
			fmt.Println("Received message")
			fmt.Println("Message Type:", messageType)
			fmt.Println(string(message))
			if err != nil {
				// controller.RemoveClient(clientId)
				conn.Close()
				log.Println(err)
				return
			}
			controller.broadcast(messageType, message)
		}
	}()
}