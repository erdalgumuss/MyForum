package models

import (
	"database/sql"
	"time"
)

type Like struct {
	ID        uint          `json:"id"`
	UserID    int           `json:"user_id"`
	PostID    sql.NullInt64 `json:"post_id"`
	CommentID sql.NullInt64 `json:"comment_id"`
	PostTitle string        `json:"post_title"`
	Action    string        `json:"action"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
