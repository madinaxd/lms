package service

import (
	"context"
	"time"

	"lms/courses_svc/client"

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

func (s *instrumentingService) GetCourse(ctx context.Context, id string) (Course, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "get").Add(1)
		s.requestLatency.With("method", "get").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GetCourse(ctx, id)
}
func (s *instrumentingService) GetCourseList(ctx context.Context, limit int, offset int) ([]Course, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "list").Add(1)
		s.requestLatency.With("method", "list").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GetCourseList(ctx, limit, offset)
}
func (s *instrumentingService) CreateCourse(ctx context.Context, name string) (Course, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "create").Add(1)
		s.requestLatency.With("method", "create").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.CreateCourse(ctx, name)
}
func (s *instrumentingService) UpdateCourse(ctx context.Context, id string, course Course) (Course, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "put").Add(1)
		s.requestLatency.With("method", "put").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.UpdateCourse(ctx, id, course)
}
func (s *instrumentingService) DeleteCourse(ctx context.Context, id string) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "delete").Add(1)
		s.requestLatency.With("method", "delete").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.DeleteCourse(ctx, id)
}
func (s *instrumentingService) GetCourseStudents(ctx context.Context, id string) ([]client.Student, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "getStudents").Add(1)
		s.requestLatency.With("method", "getStudents").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GetCourseStudents(ctx, id)
}
