package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

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
		log.Printf("task repository: get: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Status); err != nil {
			log.Printf("task repository: get: %v\n", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Close(); err != nil {
		log.Printf("task repository: get: %v\n", err)
		return nil, err
	}

	if err := rows.Err(); err != nil {
		log.Printf("task repository: get: %v\n", err)
		return nil, err
	}

	return tasks, nil
}

func (r *repository) Check() error {
	var task model.Task
	row := r.db.QueryRow("SELECT id, title, status FROM tasks LIMIT 1")

	err := row.Scan(&task.ID, &task.Title, &task.Status)
	if err != nil {
		log.Printf("task repository: check: %v\n", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	}

	return nil
}

func (r *repository) Create(task model.Task) (model.Task, error) {
	result, err := r.db.Exec("INSERT INTO tasks (title, status) VALUES (?, ?)", task.Title, task.Status)
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
	row := r.db.QueryRow("SELECT id, title, status FROM tasks WHERE id=?", id)

	err := row.Scan(&task.ID, &task.Title, &task.Status)
	if err != nil {
		log.Println("repository:", err)
		return task, err
	}

	return task, nil
}
