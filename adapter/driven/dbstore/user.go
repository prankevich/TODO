package dbstore

import (
	"TODO/adapter/driven/models"
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CreateTask(ctx context.Context, t models.Task) (models.Task, error) {
	t.ID = generateID()
	now := time.Now().UTC()
	t.CreatedAt = now
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO tasks (id, user_id, title, notes, due_at, done, created_at,updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		t.ID, t.UserID, t.Title, t.Notes, t.DueAt, t.Done, t.CreatedAt, t.UpdatedAt)
	return t, err
}
func (r *Repo) ListTasks(ctx context.Context, userID string) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.SelectContext(ctx, &tasks,
		`SELECT id, user_id, title, notes, due_at, done, created_at, updated_at
         FROM tasks WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	return tasks, err
}
func (r *Repo) CompleteTask(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE tasks SET done=true, updated_at=NOW() WHERE id=$1`, id)
	if err != nil {
		return err
	}
	return nil
}

func generateID() string {
	return "tsk_" + time.Now().Format("2006.01.02.000")
}
