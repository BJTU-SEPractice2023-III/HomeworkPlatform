package models

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

// Admin + gzofLRnA
func TestCheckPassword_Error(t *testing.T) {
	user := User{Username: "Admin", Password: "EpK6fkcbLQAihcD3:00d9effabb8ce9d0093be6d4faff72e694af28eb"}
	assert.Equal(t, false, user.CheckPassword("gzofLRnA"))
}

func TestCheckPassword_NULL(t *testing.T) {
	user := User{Username: "Admin", Password: "EpK6fkcbLQAihcD3:823e4bebf5915ba903d4bda434457e40d0dc789e"}
	assert.Equal(t, false, user.CheckPassword(""))
}

func TestCheckPassword_Right(t *testing.T) {
	user := User{Username: "Admin", Password: "EpK6fkcbLQAihcD3:823e4bebf5915ba903d4bda434457e40d0dc789e"}
	assert.Equal(t, true, user.CheckPassword("gzofLRnA"))
}

