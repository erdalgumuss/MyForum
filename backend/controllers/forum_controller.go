package controllers

import (
	"log"

	"MyForum/config"
	"MyForum/models"
)

func CreatePostWithPost(post models.Post) error {
	// Start a transaction
	tx, err := config.DB.Begin()
	if err != nil {
		log.Println("Transaction start error:", err)
		return err
	}

	// Insert into posts table
	query := `
	INSERT INTO posts (user_id, title, content, username, image_url, created_at)
	VALUES (?, ?, ?, ?, ?, ?)
	RETURNING id
	`
	var postID int
	err = tx.QueryRow(query, post.UserID, post.Title, post.Content, post.Username, post.ImageURL, post.CreatedAt.Format("2006-01-02 15:04:05")).Scan(&postID)
	if err != nil {
		log.Println("Error inserting post:", err)
		tx.Rollback()
		return err
	}

	// Insert categories into post_categories table
	categoryQuery := `INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`
	for _, categoryID := range post.CategoryIDs {
		_, err := tx.Exec(categoryQuery, postID, categoryID)
		if err != nil {
			log.Println("Error inserting post category:", err)
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Println("Transaction commit error:", err)
		return err
	}

	return nil
}
