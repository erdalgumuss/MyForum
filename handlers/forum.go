package handlers

import (
	"encoding/json"
	"net/http"

	"MyForum/models"
)

func CreateTopic(w http.ResponseWriter, r *http.Request) {
	var topic models.Topic
	err := json.NewDecoder(r.Body).Decode(&topic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := models.DB.Create(&topic).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(topic)
}

func ListTopics(w http.ResponseWriter, r *http.Request) {
	var topics []models.Topic
	if err := models.DB.Find(&topics).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(topics)
}
