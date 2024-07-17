package models

import "time"

type Post struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Categories string    `json:"categories"`
	Content    string    `json:"content"`
	Username   string    `json:"username"`
	UserID     int       `json:"user_id"`
	ImageURL   string    `json:"image_url"` // Yeni alan
	Likes      int       `json:"likes"`
	Dislikes   int       `json:"dislikes"`
	CreatedAt  time.Time `json:"created_at"`
}
