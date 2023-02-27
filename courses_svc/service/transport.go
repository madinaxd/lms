package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

var (
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHTTPHandler(e Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(errorEncoder),
	}
	// GET     /courses/                          retrieve courses list
	// GET     /courses/:id                       retrieve course by id
	// POST    /courses/                          adds another course
	// PUT     /courses/:id                       post updated course information about the course
	// DELETE  /courses/:id                       remove the given course
	// GET     /courses/:id/students               retrieve course students by course id
	r.Methods("GET").Path("/courses").Handler(httptransport.NewServer(
		e.GetCourseListEndpoint,
		decodeGetCourseListRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/courses/{id}").Handler(httptransport.NewServer(
		e.GetCourseEndpoint,
		decodeGetCourseRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/courses/").Handler(httptransport.NewServer(
		e.CreateCourseEndpoint,
		decodeCreateCourseRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/courses/{id}").Handler(httptransport.NewServer(
		e.UpdateCourseEndpoint,
		decodeUpdateCourseRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/courses/{id}").Handler(httptransport.NewServer(
		e.DeleteCourseEndpoint,
		decodeDeleteCourseRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/courses/{id}/students").Handler(httptransport.NewServer(
		e.GetCourseStudentsEndpoint,
		decodeGetCourseStudentsRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodeGetCourseListRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	limitstr := r.URL.Query().Get("Limit")
	offsetstr := r.URL.Query().Get("Offset")
	limit, _ := strconv.Atoi(limitstr)
	offset, _ := strconv.Atoi(offsetstr)
	return getCourseListRequest{Limit: limit, Offset: offset}, nil
}

func decodeGetCourseRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getCourseRequest{ID: id}, nil
}

func decodeCreateCourseRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req createCourseRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeUpdateCourseRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var course Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		return nil, err
	}
	return updateCourseRequest{
		ID:     id,
		Course: course,
	}, nil
}

func decodeDeleteCourseRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteCourseRequest{ID: id}, nil
}

func decodeGetCourseStudentsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getCourseStudentsRequest{ID: id}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(endpoint.Failer); ok && e.Failed() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		errorEncoder(ctx, e.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func err2code(err error) int {
	switch err {
	case ErrBadRouting, ErrInconsistentIDs:
		return http.StatusBadRequest
	case ErrNotFound:
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

type errorWrapper struct {
	Error string `json:"error"`
}
