package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	// .env dosyasını yükle
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return fmt.Errorf("error loading .env file: %v", err)
	}

	// Veritabanı bağlantısını oluştur
	dsn := os.Getenv("DB_PATH")
	DB, err = sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Veritabanı bağlantısını kontrol et
	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
		return fmt.Errorf("failed to ping database: %v", err)
	}

	fmt.Println("Database connected successfully")

	// Tabloları oluştur
	return createTables()
}

func createTables() error {
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		email TEXT,
		password TEXT,
		name TEXT,
		surname TEXT,
		created_at DATETIME,
		updated_at DATETIME

	);`
	_, err := DB.Exec(createUserTable)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
		return fmt.Errorf("failed to create users table: %v", err)
	}

	createPostTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		content TEXT,
		user_id INTEGER,
		username, TEXT,
		likes INTEGER DEFAULT 0,
		dislikes INTEGER DEFAULT 0,
		created_at DATETIME,
		updated_at DATETIME,
		image_url TEXT
	);`
	_, err = DB.Exec(createPostTable)
	if err != nil {
		log.Fatalf("Failed to create posts table: %v", err)
		return fmt.Errorf("failed to create posts table: %v", err)
	}

	createCommentTable := `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT,
		user_id INTEGER,
		post_id INTEGER,
		likes INTEGER DEFAULT 0,
		dislikes INTEGER DEFAULT 0,
		created_at DATETIME,
		updated_at DATETIME
	);`
	_, err = DB.Exec(createCommentTable)
	if err != nil {
		log.Fatalf("Failed to create comments table: %v", err)
		return fmt.Errorf("failed to create comments table: %v", err)
	}

	createCategoryTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);`
	_, err = DB.Exec(createCategoryTable)
	if err != nil {
		log.Fatalf("Failed to create categories table: %v", err)
		return fmt.Errorf("failed to create categories table: %v", err)
	}

	createPostCategoriesTable := `
	CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER,
		category_id INTEGER,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (category_id) REFERENCES categories(id)
	);`
	_, err = DB.Exec(createPostCategoriesTable)
	if err != nil {
		log.Fatalf("Failed to create post_categories table: %v", err)
		return fmt.Errorf("failed to create post_categories table: %v", err)
	}

	createSessionTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		token TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME,
		FOREIGN KEY (user_id) REFERENCES users (id)
	);
	`
	_, err = DB.Exec(createSessionTable)
	if err != nil {
		log.Fatalf("Failed to create sessions table: %v", err)
		return fmt.Errorf("failed to create sessions table: %v", err)
	}

	// Try to add the column if it doesn't exist
	alterTableQuery := `ALTER TABLE sessions ADD COLUMN expires_at DATETIME;`
	_, err = DB.Exec(alterTableQuery)
	if err != nil && !strings.Contains(err.Error(), "duplicate column name") {
		log.Fatalf("Failed to alter sessions table: %v", err)
		return fmt.Errorf("failed to alter sessions table: %v", err)
	}

	return nil
}
