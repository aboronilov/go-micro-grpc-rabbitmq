package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	if reqBody.UserID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	if reqBody.Pickup.Latitude == 0 || reqBody.Pickup.Longitude == 0 {
		http.Error(w, "Pickup is required", http.StatusBadRequest)
		return
	}

	if reqBody.Destination.Latitude == 0 || reqBody.Destination.Longitude == 0 {
		http.Error(w, "Destination is required", http.StatusBadRequest)
		return
	}

	log.Println("Preview trip request received:", reqBody)
	if err := writeJSON(w, http.StatusCreated, "Trip preview request received"); err != nil {
		log.Printf("Failed to write JSON response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
	return
}
