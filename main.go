package main

import (
	handlers "PrayerService/controller"
	"fmt"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"log"
)


func main()  {
	controller := handlers.NewController()
	err := godotenv.Load()
    if err != nil {
      log.Println("Error loading .env file")
    }
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	http.HandleFunc("/subscribe", controller.Subscribe)
	http.HandleFunc("/prayer", controller.PostPrayer)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Websocket server is running"))
	})
	if err := http.ListenAndServe(":" + port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
