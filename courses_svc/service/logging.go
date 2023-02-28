package service

import (
	"context"

	"courses/client"

	"github.com/go-kit/log"
)

type Middleware func(Service) Service

// LoggingMiddleware takes a logger as a dependency
// and returns a service Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) GetCourse(ctx context.Context, id string) (course Course, err error) {
	defer func() {
		mw.logger.Log("method", "GetCourse", "id", id, "course", course.Name, "err", err)
	}()
	return mw.next.GetCourse(ctx, id)
}
func (mw loggingMiddleware) GetCourseList(ctx context.Context, limit int, offset int) (courses []Course, err error) {
	defer func() {
		mw.logger.Log("method", "GetCourseList", "limit", limit, "offset", offset, "courses length", len(courses), "err", err)
	}()
	return mw.next.GetCourseList(ctx, limit, offset)
}
func (mw loggingMiddleware) CreateCourse(ctx context.Context, name string) (course Course, err error) {
	defer func() {
		mw.logger.Log("method", "CreateCourse", "course", course.Name, "err", err)
	}()
	return mw.next.CreateCourse(ctx, name)
}
func (mw loggingMiddleware) UpdateCourse(ctx context.Context, id string, course Course) (courseR Course, err error) {
	defer func() {
		mw.logger.Log("method", "UpdateCourse", "id", id, "course", courseR.Name, "err", err)
	}()
	return mw.next.UpdateCourse(ctx, id, course)
}
func (mw loggingMiddleware) DeleteCourse(ctx context.Context, id string) (err error) {
	defer func() {
		mw.logger.Log("method", "DeleteCourse", "id", id, "err", err)
	}()
	return mw.next.DeleteCourse(ctx, id)
}
func (mw loggingMiddleware) GetCourseStudents(ctx context.Context, id string) (students []client.Student, err error) {
	defer func() {
		mw.logger.Log("method", "GetCourseStudents", "id", id, "len", len(students), "err", err)
	}()
	return mw.next.GetCourseStudents(ctx, id)
}
