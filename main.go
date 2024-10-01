package main

import (
	controller "PrayerService/controller"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

func main() {
	controller := controller.GetInstance()
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	http.HandleFunc("/subscribe", controller.Subscribe)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Websocket server is running"))
	})
	log.Println("Server started on port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalln("Server error:", err)
	} 
}
