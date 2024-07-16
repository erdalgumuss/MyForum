package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type JSONNullString struct {
	sql.NullString
}

func (ns *JSONNullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}

func (ns JSONNullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	} else {
		return json.Marshal(nil)
	}
}

type User struct {
	ID         int            `json:"id"`    // Unique identifier for the user
	Email      string         `json:"email"` // Email of the user
	Name       JSONNullString `json:"name"`
	Surname    JSONNullString `json:"surname"`
	Username   JSONNullString `json:"username"` // Username of the user
	Password   string         `json:"password"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	GoogleID   JSONNullString `json:"googleid"`
	GitHubID   sql.NullInt64  `json:"githubid"`
	FacebookID JSONNullString `json:"facebookid"`
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
