package db

import (
	"context"
	"courses/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomCourse(t *testing.T) Course {
	Course, err := testQueries.CreateCourse(context.Background(), utils.RandomName())
	require.NoError(t, err)
	require.NotEmpty(t, Course)
	require.NotZero(t, Course.ID)
	require.NotZero(t, Course.CreatedAt)
	return Course
}

func TestCreateCourse(t *testing.T) {
	createRandomCourse(t)
}

func TestGetCourse(t *testing.T) {
	Course := createRandomCourse(t)
	id := Course.ID

	Course, err := testQueries.GetCourse(context.Background(), id)
	require.NoError(t, err)
	require.NotEmpty(t, Course)
	require.Equal(t, Course.ID, id)

}

func TestDeleteCourse(t *testing.T) {
	Course := createRandomCourse(t)
	id := Course.ID

	err := testQueries.DeleteCourse(context.Background(), id)
	require.NoError(t, err)
}

func TestListCourses(t *testing.T) {
	arg := ListCoursesParams{
		Limit:  0,
		Offset: 100,
	}
	_, err := testQueries.ListCourses(context.Background(), arg)
	require.NoError(t, err)

}

func TestUpdateCourse(t *testing.T) {
	Course := createRandomCourse(t)
	arg := UpdateCourseParams{
		ID:   Course.ID,
		Name: Course.Name,
	}

	CourseResult, err := testQueries.UpdateCourse(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Course)
	require.Equal(t, CourseResult.Name, arg.Name)
	require.Equal(t, CourseResult.CreatedAt, Course.CreatedAt)
}
