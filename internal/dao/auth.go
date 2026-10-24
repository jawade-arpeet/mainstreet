package dao

import "uuid"

type Account struct {
	ID              uuid.UUID `db:"id"`
	Email           string    `db:"email"`
	PasswordHash    string    `db:"password_hash"`
	IsActive        bool      `db:"is_active"`
	IsEmailVerified bool      `db:"is_email_verified"`
}
