package models

import "time"

type Like struct {
	ID        uint      `json:"id"`
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id"`
	CommentID int       `json:"comment_id"`
	PostTitle string    `json:"post_title"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
