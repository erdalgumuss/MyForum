package main

import (
	"MyForum/config"
	"MyForum/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Veritabanını başlat
	config.InitDB()

	// Yeni bir Gin router oluştur
	r := gin.Default()

	// Statik dosyaları sun
	r.Static("/static", "./static")

	// HTML şablonlarını yükle
	r.LoadHTMLGlob("templates/*")

	// Rotaları tanımla
	routes.AuthRoutes(r)
	routes.ForumRoutes(r)
	routes.ProfileRoutes(r) // Profil rotalarını ekledik

	// Sunucuyu başlat
	r.Run(":8080")
}
