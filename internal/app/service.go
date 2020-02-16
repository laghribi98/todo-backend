package app

import (
	"context"
	"errors"
)

// Service is a simple CRUD interface for user profiles.
type Service interface {
	GetTodos(ctx context.Context) ([]Todo, error)
	GetTodo(ctx context.Context, id int) (Todo, error)
	PostTodo(ctx context.Context, t Todo) error
	DeleteTodos(ctx context.Context) error
	DeleteTodo(ctx context.Context, id int) error
	PatchTodo(ctx context.Context, id int, t Todo) error
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type inmemService struct {
	repository Repository
}

func NewInmemService() Service {
	return &inmemService{}
}

func (s *inmemService) PostTodo(ctx context.Context, t Todo) error {
	return s.repository.Save(ctx, t)
}

func (s *inmemService) GetTodo(ctx context.Context, id int) (Todo, error) {
	return s.repository.Get(ctx, id)
}

func (s *inmemService) PatchTodo(ctx context.Context, id int, t Todo) error {
	return s.repository.Update(ctx, id, t)
}

func (s *inmemService) DeleteTodo(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}

func (s *inmemService) GetTodos(ctx context.Context) ([]Todo, error) {
	return s.repository.GetAll(ctx)
}

func (s *inmemService) DeleteTodos(ctx context.Context) error {
	return s.repository.DeleteAll(ctx)
}
