package models

import "database/sql"

type Comment struct {
	ID        uint         `json:"id"`
	Content   string       `json:"content"`
	PostID    uint         `json:"post_id"`
	UserID    int          `json:"user_id"`
	Username  string       `json:"username"`
	Likes     int          `json:"likes"`
	Dislikes  int          `json:"dislikes"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
