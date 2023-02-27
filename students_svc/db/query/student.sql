-- name: GetStudent :one
SELECT * FROM students
WHERE id = $1 LIMIT 1;

-- name: ListStudents :many
SELECT * FROM students
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CreateStudent :one
INSERT INTO students (
  fullname, date_of_birth, grade, phone
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteStudent :exec
DELETE FROM students
WHERE id = $1;

-- name: UpdateStudent :one
UPDATE students
  set fullname = $2,
  date_of_birth = $3,
  grade = $4,
  phone = $5
WHERE id = $1
RETURNING *;

-- name: GetStudentsByCourseID :many
SELECT students.id, fullname, date_of_birth, grade, phone, students.created_at FROM enrollments
INNER JOIN students
ON enrollments.student_id = students.id
WHERE course_id = $1
LIMIT $2
OFFSET $3;