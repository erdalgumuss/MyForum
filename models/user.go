package models

import "database/sql"

type User struct {
	ID       int64          `json:"id"`
	Username sql.NullString `json:"username"`
	Email    sql.NullString `json:"email"`
	Password sql.NullString `json:"password"`
	CreatedAt sql.NullTime   `json:"created_at"`
	UpdatedAt sql.NullTime   `json:"updated_at"`
}
