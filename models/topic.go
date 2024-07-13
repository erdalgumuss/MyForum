package models

type Topic struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Username string `json:"username"`
	Content  string `json:"content"`
	Category string `json:"category"`
	UserID   int    `json:"user_id"`
}
