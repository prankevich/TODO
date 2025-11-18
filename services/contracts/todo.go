package contracts

import (
	"TODO/adapter/driven/models"
	"context"
)

type TodoRepository interface {
	CreateTask(ctx context.Context, t models.Task) (models.Task, error)
	ListTasks(ctx context.Context, userID string) ([]models.Task, error)
	CompleteTask(ctx context.Context, id string) error
}
