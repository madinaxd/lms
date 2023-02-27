package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomEnrollment(t *testing.T) Enrollment {
	student := CreateRandomStudent(t)
	arg := CreateEnrollmentParams{
		StudentID: student.ID,
		CourseID:  1,
	}

	Enrollment, err := testQueries.CreateEnrollment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Enrollment)
	require.Equal(t, Enrollment.StudentID, arg.StudentID)
	require.NotZero(t, Enrollment.ID)
	require.NotZero(t, Enrollment.CreatedAt)
	return Enrollment
}

func TestCreateEnrollment(t *testing.T) {
	createRandomEnrollment(t)
}

func TestGetEnrollment(t *testing.T) {
	Enrollment := createRandomEnrollment(t)
	id := Enrollment.ID

	Enrollment, err := testQueries.GetEnrollment(context.Background(), id)
	require.NoError(t, err)
	require.NotEmpty(t, Enrollment)
	require.Equal(t, Enrollment.ID, id)

}

func TestDeleteEnrollment(t *testing.T) {
	Enrollment := createRandomEnrollment(t)
	id := Enrollment.ID

	err := testQueries.DeleteEnrollment(context.Background(), id)
	require.NoError(t, err)
}

func TestListEnrollments(t *testing.T) {
	arg := ListEnrollmentsParams{
		Limit:  0,
		Offset: 100,
	}
	_, err := testQueries.ListEnrollments(context.Background(), arg)
	require.NoError(t, err)

}
