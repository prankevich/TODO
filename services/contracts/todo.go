package contracts

import (
	"TODO/adapter/driven/models"
	"context"
	"time"
)

type TodoRepository interface {
	CreateTask(ctx context.Context, t models.Task) (models.Task, error)
	ListTasks(ctx context.Context, userID string) ([]models.Task, error)
	CompleteTask(ctx context.Context, id string) error
	DeleteTask(ctx context.Context, id string) error
	UpdateTask(ctx context.Context, t models.Task) (models.Task, error)
	ListDueTasks(ctx context.Context, before time.Time) ([]models.Task, error)
	ListDoneTasks(ctx context.Context, userID string) ([]models.Task, error)
}
