package repository

import (
	"database/sql"
	"errors"
	"fmt"

	errs "github.com/galiullindo/todolist-rest-api/internal/errors"
	"github.com/galiullindo/todolist-rest-api/internal/model"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Get() ([]model.Task, error) {
	rows, err := r.db.Query("SELECT id, title, status FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("task repository: get: %w", err)
	}
	defer rows.Close()
	tasks := make([]model.Task, 0)

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Status); err != nil {
			return nil, fmt.Errorf("task repository: get: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("task repository: get: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("task repository: get: %w", err)
	}

	return tasks, nil
}

func (r *repository) Check() error {
	var task model.Task
	row := r.db.QueryRow("SELECT id, title, status FROM tasks LIMIT 1")

	if err := row.Scan(&task.ID, &task.Title, &task.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("task repository: check: %w", err)
	}

	return nil
}

func (r *repository) Create(task model.Task) (model.Task, error) {
	result, err := r.db.Exec("INSERT INTO tasks (title, status) VALUES (@title, @status)",
		sql.Named("title", task.Title),
		sql.Named("status", task.Status),
	)

	if err != nil {
		return task, fmt.Errorf("task repository: check: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return task, fmt.Errorf("task repository: check: %w", err)
	}
	if rows != 1 {
		return task, fmt.Errorf("task repository: check: %w: expected 1, got %d", err, rows)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return task, fmt.Errorf("task repository: check: %w", err)
	}
	task.ID = id

	return task, nil
}

func (r *repository) GetByID(id int64) (model.Task, error) {
	var task model.Task
	row := r.db.QueryRow("SELECT id, title, status FROM tasks WHERE id=?",
		sql.Named("id", id),
	)

	if err := row.Scan(&task.ID, &task.Title, &task.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return task, errs.ErrTaskNotFound
		}

		return task, fmt.Errorf("task repository: check: %w", err)
	}

	return task, nil
}

func (r *repository) CheckExistenceByID(id int64) error {
	var task model.Task
	row := r.db.QueryRow("SELECT id, title, status FROM tasks WHERE id = @id",
		sql.Named("id", id),
	)

	if err := row.Scan(&task.ID, &task.Title, &task.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrTaskNotFound
		}

		return fmt.Errorf("task repository: check: %w", err)
	}

	return nil
}

func (r *repository) ReplaceByID(task model.Task) (model.Task, error) {
	result, err := r.db.Exec("UPDATE tasks SET title = @title, status = @status WHERE tasks.id = @id",
		sql.Named("id", task.ID),
		sql.Named("title", task.Title),
		sql.Named("status", task.Status),
	)

	if err != nil {
		return task, fmt.Errorf("task repository: replace by id: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return task, fmt.Errorf("task repository: replace by id: %w", err)
	}

	if rows == 0 {
		return task, errs.ErrTaskNotFound
	}
	if rows != 1 {
		return task, fmt.Errorf("task repository: replace by id: %w: expected 1, got %d", err, rows)
	}

	return task, nil
}

func (r *repository) DeleteByID(id int64) error {
	result, err := r.db.Exec("DELETE FROM tasks WHERE tasks.id = @id",
		sql.Named("id", id),
	)

	if err != nil {
		return fmt.Errorf("task repository: delete by id: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("task repository: delete by id: %w", err)
	}

	if rows == 0 {
		return errs.ErrTaskNotFound
	}
	if rows != 1 {
		return fmt.Errorf("task repository: delete by id: %w: expected 1, got %d", err, rows)
	}

	return nil
}
