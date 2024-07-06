package models

type Topic struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	UserID   int   `json:"user_id"`
}
