package service

import (
	"context"
	"time"

	"lms/students_svc/client"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/log"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

var (
	_ endpoint.Failer = getStudentResponse{}
	_ endpoint.Failer = getStudentListResponse{}
	_ endpoint.Failer = createStudentResponse{}
	_ endpoint.Failer = updateStudentResponse{}
	_ endpoint.Failer = deleteStudentResponse{}
)

type getStudentRequest struct {
	ID string
}
type getStudentResponse struct {
	Student Student `json:"student,omitempty"`
	Err     error   `json:"error,omitempty"`
}

func (r getStudentResponse) Failed() error { return r.Err }

type getStudentListRequest struct {
	Limit  int
	Offset int
}
type getStudentListResponse struct {
	Students []Student `json:"students,omitempty"`
	Err      error     `json:"error,omitempty"`
}

func (r getStudentListResponse) Failed() error { return r.Err }

type createStudentRequest struct {
	Fullname    string
	DateOfBirth time.Time
	Grade       int
	Phone       int
}

type createStudentResponse struct {
	Student Student `json:"student,omitempty"`
	Err     error   `json:"error,omitempty"`
}

func (r createStudentResponse) Failed() error { return r.Err }

type updateStudentRequest struct {
	ID      string
	Student Student
}
type updateStudentResponse struct {
	Student Student `json:"student,omitempty"`
	Err     error   `json:"error,omitempty"`
}

func (r updateStudentResponse) Failed() error { return r.Err }

type deleteStudentRequest struct {
	ID string
}
type deleteStudentResponse struct {
	Err error `json:"error,omitempty"`
}

type getStudentCoursesRequest struct {
	ID       string
	Instance string
}

type getStudentCoursesResponse struct {
	Courses []client.Course
	Err     error `json:"error,omitempty"`
}

type getCourseStudentsRequest struct {
	CourseID string
}
type getCourseStudentsResponse struct {
	Students []Student
	Err      error `json:"error,omitempty"`
}

func (r deleteStudentResponse) Failed() error { return r.Err }

type Endpoints struct {
	GetStudentEndpoint        endpoint.Endpoint
	GetStudentListEndpoint    endpoint.Endpoint
	CreateStudentEndpoint     endpoint.Endpoint
	UpdateStudentEndpoint     endpoint.Endpoint
	DeleteStudentEndpoint     endpoint.Endpoint
	GetStudentCoursesEndpoint endpoint.Endpoint
	GetCourseEndpoint         endpoint.Endpoint
	GetCourseStudentsEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(svc Service, logger log.Logger, duration metrics.Histogram) Endpoints {
	var rateLimitSeconds int = 1
	var GetStudentEndpoint endpoint.Endpoint
	{
		GetStudentEndpoint = MakeGetStudentEndpoint(svc)
		// Reqests are limited to 1 request per second with burst of 1 request.
		// Note, rate is defined as a time interval between requests.
		GetStudentEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), rateLimitSeconds))(GetStudentEndpoint)
		GetStudentEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(GetStudentEndpoint)
	}
	var GetStudentListEndpoint endpoint.Endpoint
	{
		GetStudentListEndpoint = MakeGetStudentListEndpoint(svc)
		GetStudentListEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), rateLimitSeconds))(GetStudentListEndpoint)
		GetStudentListEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(GetStudentListEndpoint)
	}
	var CreateStudentEndpoint endpoint.Endpoint
	{
		CreateStudentEndpoint = MakeCreateStudentEndpoint(svc)
		CreateStudentEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), rateLimitSeconds))(CreateStudentEndpoint)
		CreateStudentEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(CreateStudentEndpoint)
	}
	var UpdateStudentEndpoint endpoint.Endpoint
	{
		UpdateStudentEndpoint = MakeUpdateStudentEndpoint(svc)
		UpdateStudentEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), rateLimitSeconds))(UpdateStudentEndpoint)
		UpdateStudentEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(UpdateStudentEndpoint)
	}
	var DeleteStudentEndpoint endpoint.Endpoint
	{
		DeleteStudentEndpoint = MakeDeleteStudentEndpoint(svc)
		DeleteStudentEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), rateLimitSeconds))(DeleteStudentEndpoint)
		DeleteStudentEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(DeleteStudentEndpoint)
	}
	var GetStudentCoursesEndpoint endpoint.Endpoint
	{
		GetStudentCoursesEndpoint = MakeGetStudentCoursesEndpoint(svc)
		GetStudentCoursesEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), rateLimitSeconds))(GetStudentCoursesEndpoint)
		GetStudentCoursesEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(GetStudentCoursesEndpoint)
	}
	var GetCourseStudentsEndpoint endpoint.Endpoint
	{
		GetCourseStudentsEndpoint = MakeGetCourseStudentsEndpoint(svc)
		GetCourseStudentsEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), rateLimitSeconds))(GetCourseStudentsEndpoint)
		GetCourseStudentsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(GetCourseStudentsEndpoint)
	}
	return Endpoints{
		GetStudentEndpoint:        GetStudentEndpoint,
		GetStudentListEndpoint:    GetStudentListEndpoint,
		CreateStudentEndpoint:     CreateStudentEndpoint,
		UpdateStudentEndpoint:     UpdateStudentEndpoint,
		DeleteStudentEndpoint:     DeleteStudentEndpoint,
		GetStudentCoursesEndpoint: GetStudentCoursesEndpoint,
		GetCourseStudentsEndpoint: GetCourseStudentsEndpoint,
	}
}

func MakeGetStudentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getStudentRequest)
		res, e := s.GetStudent(ctx, req.ID)
		return getStudentResponse{Student: res, Err: e}, nil
	}
}
func MakeGetStudentListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getStudentListRequest)
		res, e := s.GetStudentList(ctx, req.Limit, req.Offset)
		return getStudentListResponse{Students: res, Err: e}, nil
	}
}

func MakeCreateStudentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createStudentRequest)
		res, e := s.CreateStudent(ctx, req.Fullname, req.DateOfBirth, req.Grade, req.Phone)
		return createStudentResponse{Student: res, Err: e}, nil
	}
}

func MakeUpdateStudentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateStudentRequest)
		res, e := s.UpdateStudent(ctx, req.ID, req.Student)
		return updateStudentResponse{Student: res, Err: e}, nil
	}
}

func MakeDeleteStudentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteStudentRequest)
		e := s.DeleteStudent(ctx, req.ID)
		return deleteStudentResponse{Err: e}, nil
	}
}

func MakeGetStudentCoursesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getStudentCoursesRequest)
		res, e := s.GetStudentCourses(ctx, req.ID)
		return getStudentCoursesResponse{Courses: res, Err: e}, nil
	}
}
func MakeGetCourseStudentsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getCourseStudentsRequest)
		res, e := s.GetCourseStudents(ctx, req.CourseID)
		return getCourseStudentsResponse{Students: res, Err: e}, nil
	}
}
