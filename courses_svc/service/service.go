package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"lms/courses_svc/client"
	db "lms/courses_svc/db/sqlc"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
)

type Service interface {
	GetCourse(ctx context.Context, id string) (Course, error)
	GetCourseList(ctx context.Context, limit int, offset int) ([]Course, error)
	CreateCourse(ctx context.Context, name string) (Course, error)
	UpdateCourse(ctx context.Context, id string, Course Course) (Course, error)
	DeleteCourse(ctx context.Context, id string) error
	GetCourseStudents(ctx context.Context, id string) ([]client.Student, error)
}

func New(r *db.Queries, studentsSvc client.StudentServiceClient, logger log.Logger, counter metrics.Counter, latency metrics.Histogram) Service {
	var svc Service
	{
		svc = NewCourseService(r, studentsSvc)
		svc = LoggingMiddleware(logger)(svc)
		svc = NewInstrumentingService(counter, latency, svc)
	}
	return svc
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
	ErrDB              = errors.New("db error")
)

func NewCourseService(r *db.Queries, studentsSvc client.StudentServiceClient) Service {
	return &CourseService{
		r:           r,
		studentsSvc: studentsSvc,
	}
}

type CourseService struct {
	r           *db.Queries
	studentsSvc client.StudentServiceClient
}

type Course struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name"`
}

func (s *CourseService) GetCourse(ctx context.Context, id string) (Course, error) {
	ID, err := strconv.Atoi(id)
	if err != nil {
		return Course{}, ErrInconsistentIDs
	}
	result, err := s.r.GetCourse(ctx, int64(ID))
	if err != nil {
		return Course{}, ErrDB
	}
	return Course{
		ID:   result.ID,
		Name: result.Name,
	}, nil
}

func (s *CourseService) GetCourseList(ctx context.Context, limit int, offset int) ([]Course, error) {
	if limit == 0 {
		limit = 100
	}
	p, err := s.r.ListCourses(ctx, db.ListCoursesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		fmt.Print(err)
		return []Course{}, ErrDB
	}
	var list []Course
	for _, result := range p {
		list = append(list, Course{
			ID:   result.ID,
			Name: result.Name,
		})
	}
	return list, nil
}

func (s *CourseService) CreateCourse(ctx context.Context, name string) (Course, error) {
	result, err := s.r.CreateCourse(ctx, name)
	if err != nil {
		return Course{}, ErrDB
	}
	return Course{
		ID:   result.ID,
		Name: result.Name,
	}, nil
}

func (s *CourseService) UpdateCourse(ctx context.Context, id string, course Course) (Course, error) {
	ID, err := strconv.Atoi(id)
	if err != nil {
		return Course{}, ErrInconsistentIDs
	}
	result, err := s.r.UpdateCourse(ctx, db.UpdateCourseParams{
		ID:   int64(ID),
		Name: course.Name,
	})
	if err != nil {
		return Course{}, ErrDB
	}
	return Course{
		ID:   result.ID,
		Name: result.Name,
	}, nil
}

func (s *CourseService) DeleteCourse(ctx context.Context, id string) error {
	ID, err := strconv.Atoi(id)
	if err != nil {
		return ErrInconsistentIDs
	}
	err = s.r.DeleteCourse(ctx, int64(ID))
	if err != nil {
		return ErrDB
	}
	return nil
}

func (s *CourseService) GetCourseStudents(ctx context.Context, id string) ([]client.Student, error) {
	res, err := s.studentsSvc.GetCourseStudents(ctx, id)

	if err != nil {
		return []client.Student{}, err
	}

	return res, nil
}
