package main

import (
	"PrayerService/controller"
	"fmt"
	"net/http"
)


func main()  {
	controller := controller.NewController()

	http.HandleFunc("/subscribe", controller.Subscribe)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}