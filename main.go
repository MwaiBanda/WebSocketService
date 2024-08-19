package main

import (
	"PrayerService/controller"
	"fmt"
	"net/http"
	"os"
)


func main()  {
	controller := controller.NewController()
	port := func() string {
		if len(os.Getenv("PORT")) > 0 {
			return os.Getenv("PORT")
		} else {
			return "8080"
		}
	}()
	http.HandleFunc("/subscribe", controller.Subscribe)
	if err := http.ListenAndServe(":" + port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
