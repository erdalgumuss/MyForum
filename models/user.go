package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID         int            `json:"id"`    // Unique identifier for the user
	Email      string         `json:"email"` // Email of the user (nullable)
	Name       sql.NullString `json:"name"`
	Surname    sql.NullString `json:"surname"`
	Username   sql.NullString `json:"username"` // Username of the user (nullable)
	Password   string         `json:"password"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	GoogleID   sql.NullString `json:"googleid"`
	GitHubID   sql.NullInt64  `json:"githubid"`
	FacebookID sql.NullString `json:"facebookid"`
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
