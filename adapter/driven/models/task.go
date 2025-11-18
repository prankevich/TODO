package models

import "time"

type Task struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	Title     string     `db:"title"`
	Notes     string     `db:"notes"`
	DueAt     *time.Time `db:"due_at"`
	Done      bool       `db:"done"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
}
