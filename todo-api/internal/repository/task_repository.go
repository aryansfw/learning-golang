package repository

import (
	"database/sql"
	"fmt"
	"todo/internal/models"

	"github.com/google/uuid"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) List(userId uuid.UUID, filters models.TaskFilter) (*[]models.Task, error) {
	var tasks []models.Task

	var query = "SELECT * FROM tasks WHERE user_id = $1"
	var args []any = []any{userId}

	var i = 2

	if filters.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", i)
		args = append(args, filters.Status)
		i++
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.Id, &task.UserId, &task.Title, &task.Description, &task.Status); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &tasks, nil
}

func (r *TaskRepository) Create(task models.Task) (*models.Task, error) {
	var resultTask models.Task
	if err := r.db.QueryRow(
		"INSERT INTO tasks VALUES ($1, $2, $3, $4, $5) "+
			"RETURNING id, user_id, title, description, status",
		task.Id, task.UserId, task.Title, task.Description, task.Status,
	).Scan(&resultTask.Id, &resultTask.UserId, &resultTask.Title,
		&resultTask.Description, &resultTask.Status); err != nil {
		return nil, err
	}

	return &resultTask, nil
}

func (r *TaskRepository) Update(task models.Task) (*models.Task, error) {
	var resultTask models.Task
	if err := r.db.QueryRow(
		"UPDATE tasks SET title = $1, description = $2, status = $3 WHERE id = $4 "+
			"RETURNING id, user_id, title, description, status",
		task.Title, task.Description, task.Status, task.Id,
	).Scan(&resultTask.Id, &resultTask.UserId, &resultTask.Title,
		&resultTask.Description, &resultTask.Status); err != nil {
		return nil, err
	}

	return &resultTask, nil
}

func (r *TaskRepository) Delete(id uuid.UUID) error {
	res, err := r.db.Exec("DELETE FROM tasks WHERE id = $1", id)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
