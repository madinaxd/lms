package service

import (
	"context"
	"time"

	"courses/client"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/log"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

var (
	_ endpoint.Failer = getCourseResponse{}
	_ endpoint.Failer = getCourseListResponse{}
	_ endpoint.Failer = createCourseResponse{}
	_ endpoint.Failer = updateCourseResponse{}
	_ endpoint.Failer = deleteCourseResponse{}
)

type getCourseRequest struct {
	ID string
}
type getCourseResponse struct {
	Course Course `json:"course,omitempty"`
	Err    error  `json:"error,omitempty"`
}

func (r getCourseResponse) Failed() error { return r.Err }

type getCourseListRequest struct {
	Limit  int
	Offset int
}
type getCourseListResponse struct {
	Courses []Course `json:"courses,omitempty"`
	Err     error    `json:"error,omitempty"`
}

func (r getCourseListResponse) Failed() error { return r.Err }

type createCourseRequest struct {
	Name string
}

type createCourseResponse struct {
	Course Course `json:"course,omitempty"`
	Err    error  `json:"error,omitempty"`
}

func (r createCourseResponse) Failed() error { return r.Err }

type updateCourseRequest struct {
	ID     string
	Course Course
}
type updateCourseResponse struct {
	Course Course `json:"course,omitempty"`
	Err    error  `json:"error,omitempty"`
}

func (r updateCourseResponse) Failed() error { return r.Err }

type deleteCourseRequest struct {
	ID string
}
type deleteCourseResponse struct {
	Err error `json:"error,omitempty"`
}

type getCourseStudentsRequest struct {
	ID       string
	Instance string
}

type getCourseStudentsResponse struct {
	Students []client.Student
	Err      error `json:"error,omitempty"`
}

func (r deleteCourseResponse) Failed() error { return r.Err }

type Endpoints struct {
	GetCourseEndpoint         endpoint.Endpoint
	GetCourseListEndpoint     endpoint.Endpoint
	CreateCourseEndpoint      endpoint.Endpoint
	UpdateCourseEndpoint      endpoint.Endpoint
	DeleteCourseEndpoint      endpoint.Endpoint
	GetCourseStudentsEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(svc Service, logger log.Logger, duration metrics.Histogram) Endpoints {
	var rateLimitSeconds int = 1
	var GetCourseEndpoint endpoint.Endpoint
	{
		GetCourseEndpoint = MakeGetCourseEndpoint(svc)
		// Reqests are limited to 1 request per second with burst of 1 request.
		// Note, rate is defined as a time interval between requests.
		GetCourseEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), rateLimitSeconds))(GetCourseEndpoint)
		GetCourseEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(GetCourseEndpoint)
	}
	var GetCourseListEndpoint endpoint.Endpoint
	{
		GetCourseListEndpoint = MakeGetCourseListEndpoint(svc)
		GetCourseListEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), rateLimitSeconds))(GetCourseListEndpoint)
		GetCourseListEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(GetCourseListEndpoint)
	}
	var CreateCourseEndpoint endpoint.Endpoint
	{
		CreateCourseEndpoint = MakeCreateCourseEndpoint(svc)
		CreateCourseEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), rateLimitSeconds))(CreateCourseEndpoint)
		CreateCourseEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(CreateCourseEndpoint)
	}
	var UpdateCourseEndpoint endpoint.Endpoint
	{
		UpdateCourseEndpoint = MakeUpdateCourseEndpoint(svc)
		UpdateCourseEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), rateLimitSeconds))(UpdateCourseEndpoint)
		UpdateCourseEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(UpdateCourseEndpoint)
	}
	var DeleteCourseEndpoint endpoint.Endpoint
	{
		DeleteCourseEndpoint = MakeDeleteCourseEndpoint(svc)
		DeleteCourseEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), rateLimitSeconds))(DeleteCourseEndpoint)
		DeleteCourseEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(DeleteCourseEndpoint)
	}
	var GetCourseStudentsEndpoint endpoint.Endpoint
	{
		GetCourseStudentsEndpoint = MakeGetCourseStudentsEndpoint(svc)
		GetCourseStudentsEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), rateLimitSeconds))(GetCourseStudentsEndpoint)
		GetCourseStudentsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(GetCourseStudentsEndpoint)
	}
	return Endpoints{
		GetCourseEndpoint:         GetCourseEndpoint,
		GetCourseListEndpoint:     GetCourseListEndpoint,
		CreateCourseEndpoint:      CreateCourseEndpoint,
		UpdateCourseEndpoint:      UpdateCourseEndpoint,
		DeleteCourseEndpoint:      DeleteCourseEndpoint,
		GetCourseStudentsEndpoint: GetCourseStudentsEndpoint,
	}
}

func MakeGetCourseEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getCourseRequest)
		res, e := s.GetCourse(ctx, req.ID)
		return getCourseResponse{Course: res, Err: e}, nil
	}
}
func MakeGetCourseListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getCourseListRequest)
		res, e := s.GetCourseList(ctx, req.Limit, req.Offset)
		return getCourseListResponse{Courses: res, Err: e}, nil
	}
}

func MakeCreateCourseEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createCourseRequest)
		res, e := s.CreateCourse(ctx, req.Name)
		return createCourseResponse{Course: res, Err: e}, nil
	}
}

func MakeUpdateCourseEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateCourseRequest)
		res, e := s.UpdateCourse(ctx, req.ID, req.Course)
		return updateCourseResponse{Course: res, Err: e}, nil
	}
}

func MakeDeleteCourseEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteCourseRequest)
		e := s.DeleteCourse(ctx, req.ID)
		return deleteCourseResponse{Err: e}, nil
	}
}

func MakeGetCourseStudentsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getCourseStudentsRequest)
		res, e := s.GetCourseStudents(ctx, req.ID)
		return getCourseStudentsResponse{Students: res, Err: e}, nil
	}
}
