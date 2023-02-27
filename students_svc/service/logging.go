package service

import (
	"context"
	"students/client"
	"time"

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

func (mw loggingMiddleware) GetStudent(ctx context.Context, id string) (student Student, err error) {
	defer func() {
		mw.logger.Log("method", "GetStudent", "id", id, "student", student, "err", err)
	}()
	return mw.next.GetStudent(ctx, id)
}
func (mw loggingMiddleware) GetStudentList(ctx context.Context, limit int, offset int) (students []Student, err error) {
	defer func() {
		mw.logger.Log("method", "GetStudentList", "limit", limit, "offset", offset, "students length", len(students), "err", err)
	}()
	return mw.next.GetStudentList(ctx, limit, offset)
}
func (mw loggingMiddleware) CreateStudent(ctx context.Context, fullname string, dateofbirth time.Time, grade int, phone int) (student Student, err error) {
	defer func() {
		mw.logger.Log("method", "CreateStudent", "student", student, "err", err)
	}()
	return mw.next.CreateStudent(ctx, fullname, dateofbirth, grade, phone)
}
func (mw loggingMiddleware) UpdateStudent(ctx context.Context, id string, student Student) (studentR Student, err error) {
	defer func() {
		mw.logger.Log("method", "UpdateStudent", "id", id, "student", studentR, "err", err)
	}()
	return mw.next.UpdateStudent(ctx, id, student)
}
func (mw loggingMiddleware) DeleteStudent(ctx context.Context, id string) (err error) {
	defer func() {
		mw.logger.Log("method", "DeleteStudent", "id", id, "err", err)
	}()
	return mw.next.DeleteStudent(ctx, id)
}
func (mw loggingMiddleware) GetStudentCourses(ctx context.Context, id string) (courses []client.Course, err error) {
	defer func() {
		mw.logger.Log("method", "GetCourses", "id", id, "len", len(courses), "err", err)
	}()
	return mw.next.GetStudentCourses(ctx, id)
}

func (mw loggingMiddleware) GetCourseStudents(ctx context.Context, id string) (students []Student, err error) {
	defer func() {
		mw.logger.Log("method", "GetCourseStudents", "id", id, "len", len(students), "err", err)
	}()
	return mw.next.GetCourseStudents(ctx, id)
}
