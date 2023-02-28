package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	db "students/db/sqlc"

	"students/client"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
)

type Service interface {
	GetStudent(ctx context.Context, id string) (Student, error)
	GetStudentList(ctx context.Context, limit int, offset int) ([]Student, error)
	CreateStudent(ctx context.Context, fullname string, dateofbirth time.Time, grade int, phone int) (Student, error)
	UpdateStudent(ctx context.Context, id string, student Student) (Student, error)
	DeleteStudent(ctx context.Context, id string) error
	GetStudentCourses(ctx context.Context, id string) ([]client.Course, error)
	GetCourseStudents(ctx context.Context, id string) ([]Student, error)
}

func New(r *db.Queries, courseSvc client.CourseServiceClient, logger log.Logger, counter metrics.Counter, latency metrics.Histogram) Service {
	var svc Service
	{
		svc = NewStudentService(r, courseSvc)
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

func NewStudentService(r *db.Queries, courseSvc client.CourseServiceClient) Service {
	return &studentService{
		r:         r,
		CourseSvc: courseSvc,
	}
}

type studentService struct {
	r         *db.Queries
	CourseSvc client.CourseServiceClient
}

type Student struct {
	ID          int64     `json:"id,omitempty"`
	Fullname    string    `json:"fullname"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Grade       int       `json:"grade,omitempty"`
	Phone       int64     `json:"phone"`
}

func (s *studentService) GetStudent(ctx context.Context, id string) (Student, error) {
	ID, err := strconv.Atoi(id)
	if err != nil {
		return Student{}, ErrInconsistentIDs
	}
	result, err := s.r.GetStudent(ctx, int64(ID))
	if err != nil {
		return Student{}, ErrDB
	}
	return Student{
		ID:          result.ID,
		Fullname:    result.Fullname,
		DateOfBirth: result.DateOfBirth,
	}, nil
}

func (s *studentService) GetStudentList(ctx context.Context, limit int, offset int) ([]Student, error) {
	if limit == 0 {
		limit = 100
	}
	p, err := s.r.ListStudents(ctx, db.ListStudentsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		fmt.Print(err)
		return []Student{}, ErrDB
	}
	var list []Student
	for _, result := range p {
		list = append(list, Student{
			ID:          result.ID,
			Fullname:    result.Fullname,
			DateOfBirth: result.DateOfBirth,
			Grade:       int(result.Grade),
			Phone:       result.Phone,
		})
	}
	return list, nil
}

func (s *studentService) CreateStudent(ctx context.Context, fullname string, dateofbirth time.Time, grade int, phone int) (Student, error) {
	result, err := s.r.CreateStudent(ctx, db.CreateStudentParams{
		Fullname:    fullname,
		DateOfBirth: dateofbirth,
		Grade:       int32(grade),
		Phone:       int64(phone),
	})
	if err != nil {
		return Student{}, ErrDB
	}
	return Student{
		ID:          result.ID,
		Fullname:    result.Fullname,
		DateOfBirth: result.DateOfBirth,
		Grade:       int(result.Grade),
		Phone:       result.Phone,
	}, nil
}

func (s *studentService) UpdateStudent(ctx context.Context, id string, student Student) (Student, error) {
	ID, err := strconv.Atoi(id)
	if err != nil {
		return Student{}, ErrInconsistentIDs
	}
	result, err := s.r.UpdateStudent(ctx, db.UpdateStudentParams{
		ID:          int64(ID),
		Fullname:    student.Fullname,
		DateOfBirth: student.DateOfBirth,
		Grade:       int32(student.Grade),
		Phone:       int64(student.Phone),
	})
	if err != nil {
		return Student{}, ErrDB
	}
	return Student{
		ID:          result.ID,
		Fullname:    result.Fullname,
		DateOfBirth: result.DateOfBirth,
		Grade:       int(result.Grade),
		Phone:       result.Phone,
	}, nil
}

func (s *studentService) DeleteStudent(ctx context.Context, id string) error {
	ID, err := strconv.Atoi(id)
	if err != nil {
		return ErrInconsistentIDs
	}
	err = s.r.DeleteStudent(ctx, int64(ID))
	if err != nil {
		return ErrDB
	}
	return nil
}

func (s *studentService) GetStudentCourses(ctx context.Context, id string) ([]client.Course, error) {
	ID, err := strconv.Atoi(id)
	if err != nil {
		return nil, ErrInconsistentIDs
	}
	enrollments, err := s.r.GetEnrollmentsByStudentID(ctx, db.GetEnrollmentsByStudentIDParams{
		StudentID: int64(ID),
		Offset:    0,
		Limit:     1000,
	})
	if err != nil {
		return nil, err
	}

	results := make(chan client.Course)
	errs := make(chan error)

	courseList := make([]client.Course, 0)

	for _, v := range enrollments {
		go func(courseID int64) {
			id := strconv.Itoa(int(courseID))
			course, err := s.CourseSvc.GetCourse(ctx, id)
			errs <- err
			results <- course
		}(v.CourseID)
	}
	for i := 0; i < len(enrollments); i++ {
		err := <-errs
		if err != nil {
			return nil, err
		}
		courseList = append(courseList, <-results)
	}

	return courseList, nil
}

func (s *studentService) GetCourseStudents(ctx context.Context, id string) ([]Student, error) {
	ID, err := strconv.Atoi(id)
	if err != nil {
		return nil, ErrInconsistentIDs
	}
	res, err := s.r.GetStudentsByCourseID(ctx, db.GetStudentsByCourseIDParams{
		CourseID: int64(ID),
		Offset:   0,
		Limit:    1000,
	})
	if err != nil {
		return nil, err
	}
	list := make([]Student, len(res))
	for _, result := range res {
		list = append(list, Student{
			ID:          result.ID,
			Fullname:    result.Fullname,
			DateOfBirth: result.DateOfBirth,
			Grade:       int(result.Grade),
			Phone:       result.Phone,
		})
	}
	return list, nil
}
