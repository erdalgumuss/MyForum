package models

type Comment struct {
	ID       uint   `json:"id"`
	Content  string `json:"content"`
	PostID   uint   `json:"post_id"`
	UserID   int   `json:"user_id"`
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
}
