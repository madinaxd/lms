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
	// GET     /students/                          retrieve students list
	// GET     /students/:id                       retrieve student by id
	// POST    /students/                          adds another student
	// PUT     /students/:id                       post updated student information about the student
	// DELETE  /students/:id                       remove the given student
	// GET     /students/:id/courses               retrieve student courses by student id
	// GET	   /courses/:id/students			   retrieve students by course id
	r.Methods("GET").Path("/students").Handler(httptransport.NewServer(
		e.GetStudentListEndpoint,
		decodeGetStudentListRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/students/{id}").Handler(httptransport.NewServer(
		e.GetStudentEndpoint,
		decodeGetStudentRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/students").Handler(httptransport.NewServer(
		e.CreateStudentEndpoint,
		decodeCreateStudentRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/students/{id}").Handler(httptransport.NewServer(
		e.UpdateStudentEndpoint,
		decodeUpdateStudentRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/students/{id}").Handler(httptransport.NewServer(
		e.DeleteStudentEndpoint,
		decodeDeleteStudentRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/students/{id}/courses").Handler(httptransport.NewServer(
		e.GetStudentCoursesEndpoint,
		decodeGetStudentCoursesRequest,
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

func decodeGetStudentListRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	limitstr := r.URL.Query().Get("Limit")
	offsetstr := r.URL.Query().Get("Offset")
	limit, _ := strconv.Atoi(limitstr)
	offset, _ := strconv.Atoi(offsetstr)
	return getStudentListRequest{Limit: limit, Offset: offset}, nil
}

func decodeGetStudentRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getStudentRequest{ID: id}, nil
}

func decodeCreateStudentRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req createStudentRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeUpdateStudentRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var student Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		return nil, err
	}
	return updateStudentRequest{
		ID:      id,
		Student: student,
	}, nil
}

func decodeDeleteStudentRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteStudentRequest{ID: id}, nil
}

func decodeGetStudentCoursesRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getStudentCoursesRequest{ID: id}, nil
}

func decodeGetCourseStudentsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getCourseStudentsRequest{CourseID: id}, nil
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
