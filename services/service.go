package services

import (
	"TODO/adapter/driven/models"
	"TODO/services/contracts"
	"context"
	"fmt"
	"strings"
	"time"
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
func (s *Service) DeleteTask(ctx context.Context, id string) error {
	return s.repo.DeleteTask(ctx, id)
}
func (s *Service) UpdateTask(ctx context.Context, t models.Task) (models.Task, error) {
	if t.ID == "" {
		return t, fmt.Errorf("task ID is required")
	}

	// Можно добавить бизнес‑правила:
	// например, запрещать обновление завершённых задач
	if t.Done {
		return t, fmt.Errorf("cannot update completed task")
	}

	// Делегируем обновление репозиторию
	updatedTask, err := s.repo.UpdateTask(ctx, t)
	if err != nil {
		return t, err
	}

	return updatedTask, nil
}

func (s *Service) ListDueTasks(ctx context.Context, before time.Time) ([]models.Task, error) {
	return s.repo.ListDueTasks(ctx, before)
}

func (s *Service) ListDoneTasks(ctx context.Context, userID string) ([]models.Task, error) {
	return s.repo.ListDoneTasks(ctx, userID)
}
