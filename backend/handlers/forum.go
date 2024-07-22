package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"MyForum/config"
	"MyForum/controllers"
	"MyForum/models"
	"MyForum/utils"

	"github.com/gin-gonic/gin"
)

func ShowForumPage(c *gin.Context) {
	c.HTML(http.StatusOK, "forum.html", gin.H{
		"Title": "Forum",
	})
}

func HandleForumPage(c *gin.Context) {
	// Handle form submission or other POST data processing here
	// For now, just render the forum page as a placeholder
	c.HTML(http.StatusOK, "forum.html", gin.H{
		"title":   "Forum",
		"message": "Post received",
	})
}

func RenderCreatePostPage(c *gin.Context) {
	c.HTML(http.StatusOK, "create_post.html", nil)
}

func GalleryPage(c *gin.Context) {
	c.HTML(http.StatusOK, "gallery.html", gin.H{
		"Title": "Galeri",
	})
}

func RulesPage(c *gin.Context) {
	c.HTML(http.StatusOK, "rules.html", gin.H{
		"Title": "Kurallar",
	})
}

func CreatePost(c *gin.Context) {
	log.Println("CreatePost function called in handlers")

	var input models.Post

	// Get user ID and username from context
	userID, ok := utils.GetUserIDFromSession(c)
	if !ok {
		log.Println("User ID not found in session")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	username, ok := c.Get("username")
	if !ok {
		log.Println("Username not found in session")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Bind form data explicitly
	input.Title = c.PostForm("title")
	categoryNames := c.PostFormArray("categories")
	input.Content = c.PostForm("content")
	input.Username = username.(string)
	input.UserID = userID
	input.CreatedAt = time.Now()

	// Log received form data
	log.Printf("Received form data: Title=%s, Content=%s, Username=%s, UserID=%d, Categories=%v\n",
		input.Title, input.Content, input.Username, input.UserID, categoryNames)

	// Handle file upload
	file, err := c.FormFile("image")
	if err == nil {
		filename := filepath.Base(file.Filename)
		filepath := filepath.Join("uploads", filename)
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			log.Println("Error saving file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "File could not be saved"})
			return
		}
		input.ImageURL = "/uploads/" + filename
		log.Printf("File uploaded successfully: %s\n", filepath)
	} else {
		log.Println("No file uploaded")
	}

	// Convert category names to IDs
	var categoryIDs []int
	for _, categoryName := range categoryNames {
		categoryName = strings.TrimSpace(categoryName)
		var categoryID int
		err := config.DB.QueryRow("SELECT id FROM categories WHERE name = ?", categoryName).Scan(&categoryID)
		if err != nil {
			log.Println("Invalid category name:", categoryName)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category name: " + categoryName})
			return
		}
		categoryIDs = append(categoryIDs, categoryID)
	}
	input.CategoryIDs = categoryIDs

	log.Printf("Form data bound successfully: %+v\n", input)

	// Call controller function
	if err := controllers.CreatePostWithPost(input); err != nil {
		log.Println("Error creating post:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create post"})
		return
	} else {
		c.Redirect(http.StatusFound, "/forum")
	}

	//c.JSON(http.StatusOK, gin.H{"message": "Post created successfully"})
}

func GetPosts(c *gin.Context) {
	var posts []models.Post

	// Query posts with sorting by created_at descending
	rows, err := config.DB.Query("SELECT id, COALESCE(username, '') AS username, title, content, user_id, likes, dislikes, created_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	defer rows.Close()

	// Iterate over rows and scan into Post structs
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Username, &post.Title, &post.Content, &post.UserID, &post.Likes, &post.Dislikes, &post.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan posts"})
			return
		}
		posts = append(posts, post)
	}

	// Handle any iteration errors
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating through posts"})
		return
	}

	// Return posts as JSON response
	c.JSON(http.StatusOK, posts)
}

func GetPost(c *gin.Context) {
	// Extract post ID from URL parameter
	id := c.Param("id")

	var post models.Post
	var categoryIDs []int
	var categoryNames []string
	var createdAt sql.NullTime

	// Query to fetch post details including username from users table
	err := config.DB.QueryRow(`
		SELECT p.id, p.title, p.content, p.likes, p.dislikes, p.user_id, p.image_url, p.created_at, u.username
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`, id).Scan(&post.ID, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &post.UserID, &post.ImageURL, &createdAt, &post.Username)

	if err != nil {
		log.Println("Error fetching post:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	formattedCreatedAt := "Unknown"
	if createdAt.Valid {
		formattedCreatedAt = createdAt.Time.Format("2006-01-02 15:04:05")
	}

	// Query to fetch comments for the post
	rows, err := config.DB.Query(`
		SELECT c.id, c.content, COALESCE(u.username, 'Unknown') AS username, c.post_id, c.likes, c.dislikes, c.created_at, c.updated_at 
		FROM comments c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at DESC
	`, id)
	if err != nil {
		log.Println("Failed to fetch comments from DB:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	defer rows.Close()

	var comments []map[string]interface{}
	for rows.Next() {
		var comment models.Comment
		var username string
		var createdAt, updatedAt sql.NullTime
		err := rows.Scan(&comment.ID, &comment.Content, &username, &comment.PostID, &comment.Likes, &comment.Dislikes, &createdAt, &updatedAt)
		if err != nil {
			log.Println("Failed to scan comment:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan comment"})
			return
		}
		comment.Username = username

		formattedCreatedAt := "Unknown"
		if createdAt.Valid {
			formattedCreatedAt = createdAt.Time.Format("2006-01-02 15:04:05")
		}

		formattedUpdatedAt := "Unknown"
		if updatedAt.Valid {
			formattedUpdatedAt = updatedAt.Time.Format("2006-01-02 15:04:05")
		}

		formattedComment := map[string]interface{}{
			"id":         comment.ID,
			"content":    comment.Content,
			"username":   comment.Username,
			"post_id":    comment.PostID,
			"likes":      comment.Likes,
			"dislikes":   comment.Dislikes,
			"created_at": formattedCreatedAt,
			"updated_at": formattedUpdatedAt,
		}
		comments = append(comments, formattedComment)
	}

	if err := rows.Err(); err != nil {
		log.Println("Failed to iterate comments:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate comments"})
		return
	}

	// Query to fetch category IDs and names
	rows, err = config.DB.Query(`
		SELECT c.id, c.name
		FROM categories c
		INNER JOIN post_categories pc ON c.id = pc.category_id
		WHERE pc.post_id = ?
	`, post.ID)
	if err != nil {
		log.Println("Error fetching categories:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching categories"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var categoryID int
		var categoryName string
		if err := rows.Scan(&categoryID, &categoryName); err != nil {
			log.Println("Error scanning category:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning categories"})
			return
		}
		categoryIDs = append(categoryIDs, categoryID)
		categoryNames = append(categoryNames, categoryName)
	}

	post.CategoryIDs = categoryIDs

	// Convert category names slice to a comma-separated string for display purposes
	categoriesString := strings.Join(categoryNames, ", ")

	c.HTML(http.StatusOK, "post.html", gin.H{
		"Post":       post,
		"Categories": categoriesString,
		"Comments":   comments,
		"CreatedAt":  formattedCreatedAt,
	})
}

func CreateComment(c *gin.Context) {
	// Get user ID from context
	userID, ok := utils.GetUserIDFromSession(c)
	if !ok {
		log.Println("User ID not found in session")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse form data
	content := c.PostForm("content")
	postIDStr := c.PostForm("post_id")

	// Convert post_id to integer
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Fetch username from the database
	var username string
	err = config.DB.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch username"})
		return
	}

	stmt, err := config.DB.Prepare("INSERT INTO comments (content, user_id, post_id, username, created_at, updated_at) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(content, userID, postID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	} else {
		c.Redirect(http.StatusSeeOther, "/posts/"+postIDStr)
	}

	//c.JSON(http.StatusOK, gin.H{"message": "Comment created successfully"})
}

func GetComments(c *gin.Context) {
	postID := c.Param("id")
	log.Println("Fetching comments for post ID:", postID)

	rows, err := config.DB.Query(`
		SELECT c.id, c.content, COALESCE(u.username, 'Unknown') AS username, c.post_id, c.likes, c.dislikes, c.created_at, c.updated_at 
		FROM comments c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at DESC
	`, postID)
	if err != nil {
		log.Println("Failed to fetch comments from DB:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	defer rows.Close()

	var comments []map[string]interface{}
	for rows.Next() {
		var comment models.Comment
		var username string
		var createdAt, updatedAt sql.NullTime
		err := rows.Scan(&comment.ID, &comment.Content, &username, &comment.PostID, &comment.Likes, &comment.Dislikes, &createdAt, &updatedAt)
		if err != nil {
			log.Println("Failed to scan comment:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan comment"})
			return
		}
		comment.Username = username

		formattedCreatedAt := "Unknown"
		if createdAt.Valid {
			formattedCreatedAt = createdAt.Time.Format("2006-01-02 15:04:05")
		}

		formattedUpdatedAt := "Unknown"
		if updatedAt.Valid {
			formattedUpdatedAt = updatedAt.Time.Format("2006-01-02 15:04:05")
		}

		formattedComment := map[string]interface{}{
			"id":         comment.ID,
			"content":    comment.Content,
			"username":   comment.Username,
			"post_id":    comment.PostID,
			"likes":      comment.Likes,
			"dislikes":   comment.Dislikes,
			"created_at": formattedCreatedAt,
			"updated_at": formattedUpdatedAt,
		}
		comments = append(comments, formattedComment)
	}

	if err := rows.Err(); err != nil {
		log.Println("Failed to iterate comments:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate comments"})
		return
	}

	log.Println("Comments fetched:", comments)
	c.JSON(http.StatusOK, comments)
}

func getCategoryIDByName(categoryName string) (int, error) {
	var categoryID int
	query := "SELECT id FROM categories WHERE name = ?"

	err := config.DB.QueryRow(query, categoryName).Scan(&categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("category '%s' not found", categoryName)
		}
		log.Printf("Error fetching category ID for '%s': %v\n", categoryName, err)
		return 0, err
	}

	return categoryID, nil
}

// LikePost handles the liking of a post
func LikePost(c *gin.Context) {
	postID := c.Param("id")
	userID, ok := utils.GetUserIDFromSession(c)
	if !ok {
		log.Println("User ID not found in session")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var action string
	err := config.DB.QueryRow(`
        SELECT action 
        FROM user_likes 
        WHERE user_id = ? AND post_id = ? AND comment_id IS NULL
    `, userID, postID).Scan(&action)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Failed to check user like status:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		log.Println("Failed to start transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if action == "like" {
		// User already liked the post
		c.JSON(http.StatusOK, gin.H{"message": "Post already liked"})
		return
	}

	if action == "dislike" {
		// User disliked the post, remove dislike
		_, err = tx.Exec("UPDATE posts SET dislikes = dislikes - 1 WHERE id = ?", postID)
		if err != nil {
			log.Println("Failed to remove dislike:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
			return
		}
	}

	_, err = tx.Exec("UPDATE posts SET likes = likes + 1 WHERE id = ?", postID)
	if err != nil {
		log.Println("Failed to like post:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	if action == "" {
		_, err = tx.Exec("INSERT INTO user_likes (user_id, post_id, action) VALUES (?, ?, 'like')", userID, postID)
	} else {
		_, err = tx.Exec("UPDATE user_likes SET action = 'like' WHERE user_id = ? AND post_id = ? AND comment_id IS NULL", userID, postID)
	}
	if err != nil {
		log.Println("Failed to record like:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	var likes, dislikes int
	err = tx.QueryRow("SELECT likes, dislikes FROM posts WHERE id = ?", postID).Scan(&likes, &dislikes)
	if err != nil {
		log.Println("Failed to fetch updated counts:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated counts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Post liked",
		"likes":    likes,
		"dislikes": dislikes,
	})
}

// DislikePost handles the disliking of a post
func DislikePost(c *gin.Context) {
	postID := c.Param("id")
	userID, ok := utils.GetUserIDFromSession(c)
	if !ok {
		log.Println("User ID not found in session")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var action string
	err := config.DB.QueryRow(`
        SELECT action 
        FROM user_likes 
        WHERE user_id = ? AND post_id = ? AND comment_id IS NULL
    `, userID, postID).Scan(&action)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Failed to check user like status:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		log.Println("Failed to start transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if action == "dislike" {
		// User already disliked the post
		c.JSON(http.StatusOK, gin.H{"message": "Post already disliked"})
		return
	}

	if action == "like" {
		// User liked the post, remove like
		_, err = tx.Exec("UPDATE posts SET likes = likes - 1 WHERE id = ?", postID)
		if err != nil {
			log.Println("Failed to remove like:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
			return
		}
	}

	_, err = tx.Exec("UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?", postID)
	if err != nil {
		log.Println("Failed to dislike post:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	if action == "" {
		_, err = tx.Exec("INSERT INTO user_likes (user_id, post_id, action) VALUES (?, ?, 'dislike')", userID, postID)
	} else {
		_, err = tx.Exec("UPDATE user_likes SET action = 'dislike' WHERE user_id = ? AND post_id = ? AND comment_id IS NULL", userID, postID)
	}
	if err != nil {
		log.Println("Failed to record dislike:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	var likes, dislikes int
	err = tx.QueryRow("SELECT likes, dislikes FROM posts WHERE id = ?", postID).Scan(&likes, &dislikes)
	if err != nil {
		log.Println("Failed to fetch updated counts:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated counts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Post disliked",
		"likes":    likes,
		"dislikes": dislikes,
	})
}

// LikeComment handles the liking of a comment
func LikeComment(c *gin.Context) {
	commentID := c.Param("id")
	userID, ok := utils.GetUserIDFromSession(c)
	if !ok {
		log.Println("User ID not found in session")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var action string
	err := config.DB.QueryRow(`
        SELECT action 
        FROM user_likes 
        WHERE user_id = ? AND comment_id = ? AND post_id IS NULL
    `, userID, commentID).Scan(&action)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Failed to check user like status:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		log.Println("Failed to start transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if action == "like" {
		// User already liked the comment
		c.JSON(http.StatusOK, gin.H{"message": "Comment already liked"})
		return
	}

	if action == "dislike" {
		// User disliked the comment, remove dislike
		_, err = tx.Exec("UPDATE comments SET dislikes = dislikes - 1 WHERE id = ?", commentID)
		if err != nil {
			log.Println("Failed to remove dislike:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
			return
		}
	}

	_, err = tx.Exec("UPDATE comments SET likes = likes + 1 WHERE id = ?", commentID)
	if err != nil {
		log.Println("Failed to like comment:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	if action == "" {
		_, err = tx.Exec("INSERT INTO user_likes (user_id, comment_id, action) VALUES (?, ?, 'like')", userID, commentID)
	} else {
		_, err = tx.Exec("UPDATE user_likes SET action = 'like' WHERE user_id = ? AND comment_id = ? AND post_id IS NULL", userID, commentID)
	}
	if err != nil {
		log.Println("Failed to record like:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	var likes, dislikes int
	err = tx.QueryRow("SELECT likes, dislikes FROM comments WHERE id = ?", commentID).Scan(&likes, &dislikes)
	if err != nil {
		log.Println("Failed to fetch updated counts:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated counts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Comment liked",
		"likes":    likes,
		"dislikes": dislikes,
	})
}

// DislikeComment handles the disliking of a comment
func DislikeComment(c *gin.Context) {
	commentID := c.Param("id")
	userID, ok := utils.GetUserIDFromSession(c)
	if !ok {
		log.Println("User ID not found in session")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var action string
	err := config.DB.QueryRow(`
        SELECT action 
        FROM user_likes 
        WHERE user_id = ? AND comment_id = ? AND post_id IS NULL
    `, userID, commentID).Scan(&action)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Failed to check user like status:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		log.Println("Failed to start transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if action == "dislike" {
		// User already disliked the comment
		c.JSON(http.StatusOK, gin.H{"message": "Comment already disliked"})
		return
	}

	if action == "like" {
		// User liked the comment, remove like
		_, err = tx.Exec("UPDATE comments SET likes = likes - 1 WHERE id = ?", commentID)
		if err != nil {
			log.Println("Failed to remove like:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
			return
		}
	}

	_, err = tx.Exec("UPDATE comments SET dislikes = dislikes + 1 WHERE id = ?", commentID)
	if err != nil {
		log.Println("Failed to dislike comment:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	if action == "" {
		_, err = tx.Exec("INSERT INTO user_likes (user_id, comment_id, action) VALUES (?, ?, 'dislike')", userID, commentID)
	} else {
		_, err = tx.Exec("UPDATE user_likes SET action = 'dislike' WHERE user_id = ? AND comment_id = ? AND post_id IS NULL", userID, commentID)
	}
	if err != nil {
		log.Println("Failed to record dislike:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	var likes, dislikes int
	err = tx.QueryRow("SELECT likes, dislikes FROM comments WHERE id = ?", commentID).Scan(&likes, &dislikes)
	if err != nil {
		log.Println("Failed to fetch updated counts:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated counts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Comment disliked",
		"likes":    likes,
		"dislikes": dislikes,
	})
}
