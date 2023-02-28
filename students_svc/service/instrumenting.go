package service

import (
	"context"
	"time"

	"students/client"

	"github.com/go-kit/kit/metrics"
)

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		next:           s,
	}
}

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Service
}

func (s *instrumentingService) GetStudent(ctx context.Context, id string) (Student, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "get").Add(1)
		s.requestLatency.With("method", "get").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GetStudent(ctx, id)
}
func (s *instrumentingService) GetStudentList(ctx context.Context, limit int, offset int) ([]Student, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "list").Add(1)
		s.requestLatency.With("method", "list").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GetStudentList(ctx, limit, offset)
}
func (s *instrumentingService) CreateStudent(ctx context.Context, fullname string, dateofbirth time.Time, grade int, phone int) (Student, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "create").Add(1)
		s.requestLatency.With("method", "create").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.CreateStudent(ctx, fullname, dateofbirth, grade, phone)
}
func (s *instrumentingService) UpdateStudent(ctx context.Context, id string, student Student) (Student, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "put").Add(1)
		s.requestLatency.With("method", "put").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.UpdateStudent(ctx, id, student)
}
func (s *instrumentingService) DeleteStudent(ctx context.Context, id string) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "delete").Add(1)
		s.requestLatency.With("method", "delete").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.DeleteStudent(ctx, id)
}
func (s *instrumentingService) GetStudentCourses(ctx context.Context, id string) ([]client.Course, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "getCourses").Add(1)
		s.requestLatency.With("method", "getCourses").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GetStudentCourses(ctx, id)
}
func (s *instrumentingService) GetCourseStudents(ctx context.Context, id string) ([]Student, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "getCourseStudents").Add(1)
		s.requestLatency.With("method", "getCourseStudents").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GetCourseStudents(ctx, id)
}
