-- name: GetCourse :one
SELECT * FROM Courses
WHERE id = $1 LIMIT 1;

-- name: ListCourses :many
SELECT * FROM Courses
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CreateCourse :one
INSERT INTO Courses (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: DeleteCourse :exec
DELETE FROM Courses
WHERE id = $1;

-- name: UpdateCourse :one
UPDATE Courses
  set name = $2
WHERE id = $1
RETURNING *;
