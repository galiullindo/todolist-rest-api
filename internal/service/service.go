package service

import (
	"github.com/galiullindo/todolist-rest-api/internal/model"
)

type Repository interface {
	GetList() ([]model.Task, error)
	CheckList() error
	Create(task model.Task) (model.Task, error)

	GetByID(id int64) (model.Task, error)
	CheckExistenceByID(id int64) error
	ReplaceByID(task model.Task) (model.Task, error)
	DeleteByID(id int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetList() ([]model.Task, error) {
	return s.repo.GetList()
}

func (s *service) GetListMetadata() error {
	return s.repo.CheckList()
}

func (s *service) Create(title string, status string) (model.Task, error) {
	newTask := model.Task{
		Title:  title,
		Status: status,
	}

	return s.repo.Create(newTask)
}

func (s *service) GetByID(id int64) (model.Task, error) {
	return s.repo.GetByID(id)
}

func (s *service) GetByIDMetadata(id int64) error {
	return s.repo.CheckExistenceByID(id)
}

func (s *service) ReplaceByID(id int64, title string, status string) (model.Task, error) {
	newTask := model.Task{
		ID:     id,
		Title:  title,
		Status: status,
	}

	return s.repo.ReplaceByID(newTask)
}

func (s *service) DeleteByID(id int64) error {
	return s.repo.DeleteByID(id)
}
