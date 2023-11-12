package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"playground/internal/app"
	"playground/internal/pkg/password"
)

func createRandomUser(t *testing.T) *app.User {
	hashedPassword, err := password.Hash(app.RandomString(6))
	require.NoError(t, err)

	arg := &app.User{
		Username:       app.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       app.RandomOwner(),
		Email:          app.RandomEmail(),
	}

	user, err := rm.User().Create(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.False(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := rm.User().Get(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)
	oldName := oldUser.FullName

	newFullName := app.RandomOwner()
	oldUser.FullName = newFullName
	updatedUser, err := rm.User().Update(context.Background(), oldUser)

	require.NoError(t, err)
	require.NotEqual(t, oldName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)
	oldEmail := oldUser.Email

	newEmail := app.RandomEmail()
	oldUser.Email = newEmail
	updatedUser, err := rm.User().Update(context.Background(), oldUser)

	require.NoError(t, err)
	require.NotEqual(t, oldEmail, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)
	oldHashedPassword := oldUser.HashedPassword
	newPassword := app.RandomString(6)
	newHashedPassword, err := password.Hash(newPassword)
	require.NoError(t, err)
	oldUser.HashedPassword = newHashedPassword

	updatedUser, err := rm.User().Update(context.Background(), oldUser)

	require.NoError(t, err)
	require.NotEqual(t, oldHashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)
	oldHashedPassword := oldUser.HashedPassword
	oldEmail := oldUser.Email
	oldFullName := oldUser.FullName

	newFullName := app.RandomOwner()
	newEmail := app.RandomEmail()
	newPassword := app.RandomString(6)
	newHashedPassword, err := password.Hash(newPassword)
	require.NoError(t, err)

	oldUser.FullName = newFullName
	oldUser.Email = newEmail
	oldUser.HashedPassword = newHashedPassword

	updatedUser, err := rm.User().Update(context.Background(), oldUser)

	require.NoError(t, err)
	require.NotEqual(t, oldHashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.NotEqual(t, oldEmail, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.NotEqual(t, oldFullName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)
}
