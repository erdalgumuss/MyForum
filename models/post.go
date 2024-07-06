package models

type Post struct {
	ID         uint       `json:"id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	UserID     uint       `json:"user_id"`
	Likes      int        `json:"likes"`
	Dislikes   int        `json:"dislikes"`
	Categories []Category `json:"categories"`
}
