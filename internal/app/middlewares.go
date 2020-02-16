package app

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) PostTodo(ctx context.Context, t Todo) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostTodo", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostTodo(ctx, t)
}

func (mw loggingMiddleware) GetTodo(ctx context.Context, id int) (p Todo, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetTodo", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetTodo(ctx, id)
}

func (mw loggingMiddleware) PatchTodo(ctx context.Context, id int, t Todo) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PatchTodo", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PatchTodo(ctx, id, t)
}

func (mw loggingMiddleware) DeleteTodo(ctx context.Context, id int) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteTodo", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteTodo(ctx, id)
}

func (mw loggingMiddleware) GetTodos(ctx context.Context) (addresses []Todo, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetTodos", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetTodos(ctx)
}

func (mw loggingMiddleware) DeleteTodos(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteTodos", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteTodos(ctx)
}
