package handlers

import (
	"encoding/json"
	"net/http"

	"MyForum/models"

	"github.com/gorilla/mux"
)

func CreateTopic(w http.ResponseWriter, r *http.Request) {
	var topic models.Topic
	err := json.NewDecoder(r.Body).Decode(&topic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Örneğin, oturum açmış kullanıcının kimliğini almak için bir yöntem kullanabilirsiniz
	// userID := utils.GetUserIDFromContext(r.Context())

	// Şu anda sabit olarak örneğin 1 numaralı kullanıcıyı kullanıyoruz
	topic.UserID = 1

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

func GetTopic(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	topicID := params["id"]

	var topic models.Topic
	if err := models.DB.First(&topic, topicID).Error; err != nil {
		http.Error(w, "Konu bulunamadı", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(topic)
}

// Yeni yorum ekleme
func CreateComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	topicID := params["topic_id"]

	var comment models.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	comment.TopicID = topicID

	if err := models.DB.Create(&comment).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// Yorum güncelleme veya silme
func UpdateComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	topicID := params["topic_id"]
	commentID := params["comment_id"]

	var comment models.Comment
	if err := models.DB.First(&comment, commentID).Error; err != nil {
		http.Error(w, "Yorum bulunamadı", http.StatusNotFound)
		return
	}

	// Yorum güncelleme işlemleri burada yapılabilir
	// Örneğin, comment.Body = updatedBody gibi güncelleme işlemleri

	if err := models.DB.Save(&comment).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comment)
}
