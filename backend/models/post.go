package models

import "time"

type Post struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	CategoryIDs []int  `json:"category_ids"`
	//Categories  []string  `json:"categories"`
	Content   string    `json:"content"`
	Username  string    `json:"username"`
	UserID    int       `json:"user_id"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	ImageURL  string    `json:"image_url"` // Yeni alan
	CreatedAt time.Time `json:"created_at"`
}
