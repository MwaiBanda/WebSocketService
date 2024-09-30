package controller

import (
	"PrayerService/model"
	"encoding/json"
	"net/http"
)


func (controller *Controller) PostPrayer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == http.MethodPost {
		prayer := model.Prayer{}
		err := json.NewDecoder(r.Body).Decode(&prayer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		prayer.ID = controller.prayers[len(controller.prayers)-1].ID + 1
		controller.prayers = append(controller.prayers, prayer)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(prayer)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
