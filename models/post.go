package models

import "time"

type Post struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Categories string    `json:"categories"`
	Username   string    `json:"username"`
	Content    string    `json:"content"`
	UserID     int       `json:"user_id"`
	ImageURL   string    `json:"image_url"` // Yeni alan
	Likes      int       `json:"likes"`
	Dislikes   int       `json:"dislikes"`
	CreatedAt  time.Time `json:"created_at"`
}
