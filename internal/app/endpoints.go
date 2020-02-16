package app

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects all of the endpoints that compose a profile service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them
// into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	GetTodosEndpoint    endpoint.Endpoint
	GetTodoEndpoint     endpoint.Endpoint
	PostTodoEndpoint    endpoint.Endpoint
	DeleteTodosEndpoint endpoint.Endpoint
	DeleteTodoEndpoint  endpoint.Endpoint
	PatchTodoEndpoint   endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service. Useful in a profilesvc
// server.
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostTodoEndpoint:    MakePostTodoEndpoint(s),
		GetTodoEndpoint:     MakeGetTodoEndpoint(s),
		GetTodosEndpoint:    MakeGetTodosEndpoint(s),
		PatchTodoEndpoint:   MakePatchTodoEndpoint(s),
		DeleteTodoEndpoint:  MakeDeleteTodoEndpoint(s),
		DeleteTodosEndpoint: MakeDeleteTodosEndpoint(s),
	}
}

// MakePostTodoEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePostTodoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postTodoRequest)
		e := s.PostTodo(ctx, req.Todo)
		return postTodoResponse{Err: e}, nil
	}
}

// MakeGetTodoEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetTodoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getTodoRequest)
		p, e := s.GetTodo(ctx, req.ID)
		return getTodoResponse{Profile: p, Err: e}, nil
	}
}

// MakeGetTodosEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetTodosEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		t, e := s.GetTodos(ctx)
		return getTodosResponse{Err: e, Todos: t}, nil
	}
}

// MakePatchTodoEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePatchTodoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(patchTodoRequest)
		e := s.PatchTodo(ctx, req.ID, req.Todo)
		return patchTodoResponse{Err: e}, nil
	}
}

// MakeDeleteTodoEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeDeleteTodoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteTodoRequest)
		e := s.DeleteTodo(ctx, req.ID)
		return deleteTodoResponse{Err: e}, nil
	}
}

// MakeDeleteTodosEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeDeleteTodosEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		e := s.DeleteTodos(ctx)
		return deleteTodosResponse{Err: e}, nil
	}
}

// We have two options to return errors from the business logic.
//
// We could return the error via the endpoint itself. That makes certain things
// a little bit easier, like providing non-200 HTTP responses to the client. But
// Go kit assumes that endpoint errors are (or may be treated as)
// transport-domain errors. For example, an endpoint error will count against a
// circuit breaker error count.
//
// Therefore, it's often better to return service (business logic) errors in the
// response object. This means we have to do a bit more work in the HTTP
// response encoder to detect e.g. a not-found error and provide a proper HTTP
// status code. That work is done with the errorer interface, in transport.go.
// Response types that may contain business-logic errors implement that
// interface.

type postTodoRequest struct {
	Todo Todo
}

type postTodoResponse struct {
	Err error `json:"err,omitempty"`
}

func (r postTodoResponse) error() error { return r.Err }

type getTodoRequest struct {
	ID int
}

type getTodoResponse struct {
	Profile Todo  `json:"profile,omitempty"`
	Err     error `json:"err,omitempty"`
}

func (r getTodoResponse) error() error { return r.Err }

type getTodosRequest struct {
	ID      string
	Profile Todo
}

type getTodosResponse struct {
	Err   error `json:"err,omitempty"`
	Todos []Todo
}

func (r getTodosResponse) error() error { return nil }

type patchTodoRequest struct {
	ID   int
	Todo Todo
}

type patchTodoResponse struct {
	Err error `json:"err,omitempty"`
}

func (r patchTodoResponse) error() error { return r.Err }

type deleteTodoRequest struct {
	ID int
}

type deleteTodoResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteTodoResponse) error() error { return r.Err }

type deleteTodosRequest struct {
	ProfileID string
	AddressID string
}

type deleteTodosResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteTodosResponse) error() error { return r.Err }
