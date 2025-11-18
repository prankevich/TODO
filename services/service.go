package services

import (
	"TODO/adapter/driven/models"
	"TODO/services/contracts"
	"context"
	"fmt"
	"strings"
)

type Service struct {
	repo contracts.TodoRepository
}

func NewService(repo contracts.TodoRepository) *Service {
	return &Service{repo: repo}
}

// CreateTask принимает готовую  Task
func (s *Service) CreateTask(ctx context.Context, t models.Task) (models.Task, error) {
	t.Title = strings.TrimSpace(t.Title)
	if t.Title == "" {
		return models.Task{}, fmt.Errorf("Empty title")
	}
	return s.repo.CreateTask(ctx, t)
}

func (s *Service) ListTasks(ctx context.Context, userID string) ([]models.Task, error) {
	return s.repo.ListTasks(ctx, userID)
}

func (s *Service) CompleteTask(ctx context.Context, taskID string) error {
	return s.repo.CompleteTask(ctx, taskID)
}
