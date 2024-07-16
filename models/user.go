package models

import (
	"time"
)

type User struct {
	ID         int       `json:"id"`    // Unique identifier for the user
	Email      string    `json:"email"` // Email of the user (nullable)
	Name       string    `json:"name"`
	Surname    string    `json:"surname"`
	Username   string    `json:"username"` // Username of the user (nullable)
	Password   string    `json:"password"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	google_id  string    `json:"google_id"`
	GitHubID   int       `json:"githubid"`
	FacebookID string    `json:"facebookid"`
}

type GoogleUserInfo struct {
    ID    string
    Name  string
    Email string
}

type GitHubUserInfo struct {
    ID    int
    Login string
    Email string
}

type FacebookUserInfo struct {
    ID    string
    Name  string
    Email string
}