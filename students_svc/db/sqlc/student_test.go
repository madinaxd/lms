package db

import (
	"context"
	"students/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateRandomStudent(t *testing.T) Student {
	arg := CreateStudentParams{
		Fullname:    utils.RandomName(),
		DateOfBirth: utils.RandomBirthDate(),
		Grade:       utils.RandomGrade(),
		Phone:       utils.RandomPhone(),
	}

	Student, err := testQueries.CreateStudent(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Student)
	require.Equal(t, Student.Fullname, arg.Fullname)
	// require.Equal(t, Student.DateOfBirth, arg.DateOfBirth)
	require.Equal(t, Student.Grade, arg.Grade)
	require.Equal(t, Student.Phone, arg.Phone)
	require.NotZero(t, Student.ID)
	require.NotZero(t, Student.CreatedAt)
	return Student
}

func TestCreateStudent(t *testing.T) {
	CreateRandomStudent(t)
}

func TestGetStudent(t *testing.T) {
	Student := CreateRandomStudent(t)
	id := Student.ID

	Student, err := testQueries.GetStudent(context.Background(), id)
	require.NoError(t, err)
	require.NotEmpty(t, Student)
	require.Equal(t, Student.ID, id)

}

func TestDeleteStudent(t *testing.T) {
	Student := CreateRandomStudent(t)
	id := Student.ID

	err := testQueries.DeleteStudent(context.Background(), id)
	require.NoError(t, err)
}

func TestListStudents(t *testing.T) {
	arg := ListStudentsParams{
		Limit:  0,
		Offset: 100,
	}
	_, err := testQueries.ListStudents(context.Background(), arg)
	require.NoError(t, err)

}

func TestUpdateStudent(t *testing.T) {
	Student := CreateRandomStudent(t)
	arg := UpdateStudentParams{
		ID:          Student.ID,
		Fullname:    Student.Fullname,
		DateOfBirth: Student.DateOfBirth,
		Grade:       utils.RandomGrade(),
		Phone:       utils.RandomPhone(),
	}

	StudentResult, err := testQueries.UpdateStudent(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Student)
	require.Equal(t, StudentResult.Fullname, arg.Fullname)
	// require.Equal(t, StudentResult.DateOfBirth, arg.DateOfBirth)
	require.Equal(t, StudentResult.Grade, arg.Grade)
	require.Equal(t, StudentResult.Phone, arg.Phone)
	require.Equal(t, StudentResult.CreatedAt, Student.CreatedAt)
}
