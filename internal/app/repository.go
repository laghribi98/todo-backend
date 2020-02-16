package app

import "context"

type Repository interface {
	GetAll(ctx context.Context) ([]Todo, error)
	Get(ctx context.Context, id int) (Todo, error)
	Save(ctx context.Context, t Todo) error
	DeleteAll(ctx context.Context) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, t Todo) error
}
