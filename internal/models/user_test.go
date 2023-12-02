package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	deleteData()
	assert := assert.New(t)

	user, err := CreateUser("username", "hashed_pwd;salt")
	assert.Nil(err)
	res, err := GetUserByID(user.ID)
	assert.Nil(err)
	assert.Equal(*user, res)
}

func TestDeleteUserById(t *testing.T) {
	deleteData()
	assert := assert.New(t)

	user, _ := CreateUser("username", "hashed_pwd;salt")

	res, _ := GetUserByID(user.ID)
	assert.Equal(*user, res)
	err := DeleteUserById(user.ID)
	assert.Nil(err)
	_, err = GetUserByID(user.ID)
	assert.Error(err)
}

func TestGetUsers(t *testing.T) {
	assert := assert.New(t)

	user1, _ := CreateUser("user1", "hashed_pwd;salt")
	user2, _ := CreateUser("user2", "hashed_pwd;salt")
	user3, _ := CreateUser("user3", "hashed_pwd;salt")

	users, err := GetUsers()
	assert.Nil(err)
	assert.Equal([]User{*user1, *user2, *user3}, users)

	DeleteUserById(user2.ID)

	users, err = GetUsers()
	assert.Nil(err)
	assert.Equal([]User{*user1, *user3}, users)
}