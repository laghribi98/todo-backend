package app

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
// Useful in a profilesvc server.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST    /profiles/                          adds another profile
	// GET     /profiles/:id                       retrieves the given profile by id
	// PUT     /profiles/:id                       post updated profile information about the profile
	// PATCH   /profiles/:id                       partial updated profile information
	// DELETE  /profiles/:id                       remove the given profile
	// GET     /profiles/:id/addresses/            retrieve addresses associated with the profile
	// GET     /profiles/:id/addresses/:addressID  retrieve a particular profile address
	// POST    /profiles/:id/addresses/            add a new address
	// DELETE  /profiles/:id/addresses/:addressID  remove an address

	r.Methods("GET").Path("/todos").Handler(httptransport.NewServer(
		e.GetTodosEndpoint,
		decodeGetTodosRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/todos/{id}").Handler(httptransport.NewServer(
		e.GetTodoEndpoint,
		decodeGetTodoRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/todos").Handler(httptransport.NewServer(
		e.PostTodoEndpoint,
		decodePostTodoRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/todos").Handler(httptransport.NewServer(
		e.DeleteTodosEndpoint,
		decodeDeleteTodosRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/todos/{id}").Handler(httptransport.NewServer(
		e.DeleteTodoEndpoint,
		decodeDeleteTodoRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PATCH").Path("/todos/{id}").Handler(httptransport.NewServer(
		e.PatchTodoEndpoint,
		decodePatchTodoRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodePostTodoRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postTodoRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Todo); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetTodosRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return nil, nil
}

func decodeGetTodoRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	idAsInt, err := strconv.Atoi(id)
	if !ok {
		return nil, err //TODO wrap
	}

	return getTodoRequest{ID: idAsInt}, nil
}

func decodePatchTodoRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	var t Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, err
	}

	idAsInt, err := strconv.Atoi(id)
	if !ok {
		return nil, err //TODO wrap
	}

	return patchTodoRequest{
		ID:   idAsInt,
		Todo: t,
	}, nil
}

func decodeDeleteTodosRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return nil, nil
}

func decodeDeleteTodoRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	idAsInt, err := strconv.Atoi(id)
	if !ok {
		return nil, err //TODO wrap
	}

	return deleteTodoRequest{ID: idAsInt}, nil
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
