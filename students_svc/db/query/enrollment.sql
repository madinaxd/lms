-- name: GetEnrollment :one
SELECT * FROM enrollments
WHERE id = $1 LIMIT 1;

-- name: ListEnrollments :many
SELECT * FROM enrollments
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: GetEnrollmentsByStudentID :many
SELECT * FROM enrollments
WHERE student_id = $1 
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: CreateEnrollment :one
INSERT INTO enrollments (
  student_id, course_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteEnrollment :exec
DELETE FROM enrollments
WHERE id = $1;
