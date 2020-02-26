package app

import (
	"context"
	"errors"
	"fmt"
)

// Service is a simple CRUD interface for user profiles.
type Service interface {
	GetTodos(ctx context.Context) ([]Todo, error)
	GetTodo(ctx context.Context, id int) (Todo, error)
	InsertTodo(ctx context.Context, t Todo) (Todo, error)
	DeleteTodos(ctx context.Context) error
	DeleteTodo(ctx context.Context, id int) error
	UpdateTodo(ctx context.Context, id int, t Todo) (Todo, error)
	Clear() error
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type serviceImpl struct {
	repository Repository
	cfg        *Config
}

type Todo struct {
	Id        *int    `json:"id"`
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
	Order     *int    `json:"order"`
	URL       string  `json:"url"`
}

func NewService(repository Repository, cfg *Config) Service {
	return &serviceImpl{repository: repository, cfg: cfg}
}

var blank = Todo{}

func (s *serviceImpl) InsertTodo(ctx context.Context, t Todo) (Todo, error) {
	if t.Completed == nil {
		var b bool
		t.Completed = &b
	}

	if t.Order == nil {
		var i int
		t.Order = &i
	}

	if t.Title == nil {
		var s string
		t.Title = &s
	}

	todo, err := s.repository.Save(t)
	if err != nil {
		return blank, err
	}

	return s.addURL(todo), nil
}

func (s *serviceImpl) GetTodo(_ context.Context, id int) (Todo, error) {
	todo, err := s.repository.Get(id)
	if err != nil {
		return blank, err
	}

	return s.addURL(todo), nil
}

func (s *serviceImpl) UpdateTodo(ctx context.Context, id int, t Todo) (Todo, error) {
	todo, err := s.repository.Update(id, t)
	if err != nil {
		return blank, err
	}

	return s.addURL(todo), nil
}

func (s *serviceImpl) DeleteTodo(_ context.Context, id int) error {
	return s.repository.Delete(id)
}

func (s *serviceImpl) GetTodos(_ context.Context) ([]Todo, error) {
	todos, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	if len(todos) == 0 {
		todos = []Todo{}
	}

	for i := range todos {
		todos[i] = s.addURL(todos[i])
	}

	return todos, err
}

func (s *serviceImpl) DeleteTodos(_ context.Context) error {
	return s.repository.DeleteAll()
}

func (s *serviceImpl) addURL(todo Todo) Todo {
	id := *todo.Id

	todo.URL = fmt.Sprintf("%s/%d", s.cfg.Url, id)

	return todo
}

func (s *serviceImpl) Clear() error {
	return s.repository.Drop()
}
