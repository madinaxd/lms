package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandomString(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomName() string {
	return RandomString(RandomInt(3, 8))
}

func RandomGrade() int32 {
	return int32(RandomInt(0, 11))
}

func RandomBirthDate() time.Time {
	return time.Date(RandomInt(2006, 2017), time.Month(RandomInt(1, 12)), RandomInt(1, 30), 0, 0, 0, 0, time.UTC)
}

func RandomPhone() int64 {
	return int64(RandomInt(87010000000, 87789999999))
}
