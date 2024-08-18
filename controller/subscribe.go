package controller

import (
	"fmt"
	"log"
	"net/http"
)

func (controller *Controller) Subscribe(w http.ResponseWriter, r *http.Request) {
	conn, err := controller.Upgrader.Upgrade(w, r, nil)
	fmt.Println("New connection")
	fmt.Println(conn.RemoteAddr().String())
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{controller: controller, conn: conn, send: func(messageType int, message []byte) {
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println(err)
			return
		}
	}}
	controller.AddClient(*client)
	go func()  {
		defer func() {
			fmt.Println("Lost connection")
			fmt.Println(conn.RemoteAddr().String())
			conn.Close()
		}()
		for {
			messageType, message, err := conn.ReadMessage()
			fmt.Println("Received message")
			fmt.Println(messageType)
			fmt.Println(string(message))
			if err != nil {
				log.Println(err)
				return
			}
			controller.broadcast(messageType, message)
		}
	}()

}