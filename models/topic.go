package models

type Topic struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	UserID   uint   `json:"user_id"`
}
