package routes

import (
	"MyForum/handlers"
	"net/http"
)


func SetupForumRoutes() {
    http.HandleFunc("/api/topics", handlers.CreateTopic)
    http.HandleFunc("/api/topics", handlers.ListTopics) // Yeni eklenen endpoint
    // Diğer forum işlemleri için route'ları buraya ekleyebilirsiniz
}