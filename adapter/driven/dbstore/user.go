package dbstore

import (
	"TODO/adapter/driven/models"
	"context"
	"fmt"
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
	now := time.Now().UTC()
	t.CreatedAt = now
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO tasks (user_id, title, notes, due_at, done, created_at,updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		t.UserID, t.Title, t.Notes, t.DueAt, t.Done, t.CreatedAt, t.UpdatedAt)
	return t, err
}
func (r *Repo) ListTasks(ctx context.Context, userID string) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.SelectContext(ctx, &tasks,
		`SELECT id, user_id, title, notes, due_at, done, created_at, updated_at
         FROM tasks WHERE user_id=$1 and done=false ORDER BY created_at DESC`, userID)
	return tasks, err
}
func (r *Repo) CompleteTask(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE tasks SET done=true, updated_at=NOW() WHERE id=$1`, id)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) UpdateTask(ctx context.Context, t models.Task) (models.Task, error) {
	t.UpdatedAt = time.Now().UTC()
	_, err := r.db.NamedExecContext(ctx, `
	UPDATE tasks
	SET
		title = COALESCE(NULLIF(:title, ''), title),
		notes = COALESCE(NULLIF(:notes, ''), notes),
		due_at = COALESCE(:due_at, due_at),
		done = COALESCE(:done, done),
		updated_at = :updated_at
	WHERE id = :id
`, t)
	if err != nil {
		return t, err
	}
	fmt.Println(t)
	return t, nil
}
func (r *Repo) DeleteTask(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id=$1`, id)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) ListDueTasks(ctx context.Context, before time.Time) ([]models.Task, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, title, notes, due_at, done
		FROM tasks
		WHERE due_at <= $1 AND done = false
	`, before)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Notes, &t.DueAt, &t.Done); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *Repo) ListDoneTasks(ctx context.Context, userID string) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.SelectContext(ctx, &tasks,
		`SELECT id, user_id, title, notes, due_at, done, created_at, updated_at
         FROM tasks WHERE user_id=$1 and done = true  ORDER BY created_at DESC`, userID)
	return tasks, err
}
func (r *Repo) Login(username, password string) (*models.User, error) {
	var u models.User
	err := r.db.Get(&u, "SELECT id, username, password FROM users WHERE username=$1 AND password=$2", username, password)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	return &u, nil
}
func (r *Repo) Register(user models.User) error {
	_, err := r.db.Exec(
		"INSERT INTO users (username, password) VALUES ($1, $2)",
		user.Username, user.Password,
	)
	return err
}
func (r *Repo) Stats() (map[string]interface{}, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}
	return map[string]interface{}{
		"users_count": count,
	}, nil
}
