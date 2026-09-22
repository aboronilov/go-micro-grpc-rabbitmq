package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
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

	jsonBody, _ := json.Marshal(reqBody)
	reader := bytes.NewReader(jsonBody)

	resp, err := http.Post("http://trip-service:8083/preview", "application/json", reader)
	if err != nil {
		log.Print(err)
		return
	}

	defer resp.Body.Close()

	var respBody any
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		http.Error(w, "failed to parse JSON data from trip service", http.StatusBadRequest)
		return
	}

	response := contracts.APIResponse{Data: respBody}

	writeJSON(w, http.StatusCreated, response)
}
