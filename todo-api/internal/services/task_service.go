package services

import (
	"todo/internal/models"
	"todo/internal/repository"

	"github.com/google/uuid"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) ListTasks(userId uuid.UUID) (*[]models.Task, error) {
	return s.repo.List(userId)
}

func (s *TaskService) CreateTask(
	title string,
	description string,
	status models.Status,
	userId uuid.UUID,
) (*models.Task, error) {
	task := models.Task{
		Id:          uuid.New(),
		UserId:      userId,
		Title:       title,
		Description: description,
		Status:      status,
	}

	return s.repo.Create(task)
}

func (s *TaskService) UpdateTask(id uuid.UUID, title string, description string, status models.Status) (*models.Task, error) {
	task := models.Task{
		Id:          id,
		Title:       title,
		Description: description,
		Status:      status,
	}

	return s.repo.Update(task)
}

func (s *TaskService) DeleteTask(id uuid.UUID) error {
	return s.repo.Delete(id)
}
