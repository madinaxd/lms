// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"time"
)

type Course struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}
