package database_test

import (
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/admin/services/database"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/database/repository/rdbms"
	"powerdns-auth-proxy/domain/shared/model/management"
	"powerdns-auth-proxy/domain/shared/setup"
	"powerdns-auth-proxy/domain/test"
	"testing"
)

func TestListUsers(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserService(registry, true, false)
	if err != nil {
		t.Error(err)
		return
	}

	users, err := service.ListUsers()
	if !assert.NoError(t, err, "failed to list users") {
		return
	}

	assert.NotEmpty(t, users, "expected users to not be empty")
	assert.Equal(t, 1, len(users), "expected one user")
	assert.Equal(t, test.UserUsername, users[0].Username, "expected username to match")
}

func TestListUsersComplete(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserService(registry, true, false)
	if err != nil {
		t.Error(err)
		return
	}

	users, err := service.ListUsersComplete()
	if !assert.NoError(t, err, "failed to list users") {
		return
	}

	assert.NotEmpty(t, users, "expected users to not be empty")
	assert.Equal(t, 1, len(users), "expected one user")
	assert.Equal(t, test.UserUsername, users[0].Username, "expected username to match")
}

func TestCheckUserIsNotAdmin(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserService(registry, true, false)
	if err != nil {
		t.Error(err)
		return
	}

	isAdmin, err := service.CheckIsAdminUser(registry.GetUint("userId"))
	if !assert.NoError(t, err, "failed to check user is admin") {
		return
	}

	assert.False(t, isAdmin, "expected user to not be admin")
}

func TestCheckUserIsAdmin(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserService(registry, true, true)
	if err != nil {
		t.Error(err)
		return
	}

	isAdmin, err := service.CheckIsAdminUser(registry.GetUint("userId"))
	if !assert.NoError(t, err, "failed to check user is admin") {
		return
	}

	assert.True(t, isAdmin, "expected user to be admin")
}

func TestCreateUser(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserService(registry, false, false)
	if err != nil {
		t.Error(err)
		return
	}

	payload := &management.UserCreatePayload{
		Username: test.UserUsername,
		Password: test.UserPassword,
	}

	userId, err := service.CreateUser(payload)
	if !assert.NoError(t, err, "failed to create user") {
		return
	}

	assert.NotNil(t, userId, "expected user to not be nil")

	dbUser, err := service.GetUserById(userId)
	if !assert.NoError(t, err, "failed to get user") {
		return
	}

	assert.NotNil(t, dbUser, "expected user to not be nil")
	assert.Equal(t, test.UserUsername, dbUser.Username, "expected username to match")
}

func TestGetUserById(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserService(registry, true, false)
	if err != nil {
		t.Error(err)
		return
	}

	dbUser, err := service.GetUserById(registry.GetUint("userId"))
	if !assert.NoError(t, err, "failed to get user") {
		return
	}

	assert.NotNil(t, dbUser, "expected user to not be nil")
	assert.Equal(t, test.UserUsername, dbUser.Username, "expected username to match")
}

func TestGetUserByUsername(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserService(registry, true, false)
	if err != nil {
		t.Error(err)
		return
	}

	dbUser, err := service.GetUser(test.UserUsername)
	if !assert.NoError(t, err, "failed to get user") {
		return
	}

	assert.NotNil(t, dbUser, "expected user to not be nil")
	assert.Equal(t, test.UserUsername, dbUser.Username, "expected username to match")
}

func TestDeleteUser(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserService(registry, true, false)
	if err != nil {
		t.Error(err)
		return
	}

	err = service.DeleteUser(test.UserUsername)
	if !assert.NoError(t, err, "failed to delete user") {
		return
	}

	dbUser, err := service.GetUser(test.UserUsername)
	if !assert.Error(t, err, "expected user to not exist") {
		return
	}

	assert.Nil(t, dbUser, "expected user to be nil")
}

func TestDeleteUserById(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserService(registry, true, false)
	if err != nil {
		t.Error(err)
		return
	}

	err = service.DeleteUserById(registry.GetUint("userId"))
	if !assert.NoError(t, err, "failed to delete user") {
		return
	}

	dbUser, err := service.GetUserById(registry.GetUint("userId"))
	if !assert.Error(t, err, "expected user to not exist") {
		return
	}

	assert.Nil(t, dbUser, "expected user to be nil")
}

func TestUpdateUserPasswordById(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserService(registry, true, false)
	if err != nil {
		t.Error(err)
		return
	}

	err = service.UpdateUserPasswordById(registry.GetUint("userId"), "newpassword")
	if !assert.NoError(t, err, "failed to update user password") {
		return
	}

	dbUser, err := service.GetUserById(registry.GetUint("userId"))
	if !assert.NoError(t, err, "failed to get user") {
		return
	}

	assert.NotNil(t, dbUser, "expected user to not be nil")
}

func TestProvideUserServiceFailsOnMissingUserRepository(t *testing.T) {
	container := &basics.InjectionContainer{}

	service, err := database.ProvideUserService(container)
	assert.Error(t, err, "provide user service should return an error")
	assert.Equal(t, "could not provide user service: user repository could not be resolved", err.Error())
	assert.Nil(t, service, "provide user service should return nil")
}

func TestProvideUserServiceSucceeds(t *testing.T) {
	container := &basics.InjectionContainer{}
	container.UserRepository = &rdbms.UserRepository{}

	service, err := database.ProvideUserService(container)
	assert.NoError(t, err, "provide user service should not return an error")
	assert.NotNil(t, service, "provide user service should return a service")
}

func getUserService(registry *test.Registry, withData bool, isAdmin bool) (*database.UserService, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockRepositories: true})
	if err != nil {
		return nil, err
	}

	service, err := database.ProvideUserService(container)
	if err != nil {
		return nil, err
	}

	if withData {
		err = test.RepositoryCreateTestUser(container, registry)
		if err != nil {
			return nil, err
		}

		if isAdmin {
			err = test.RepositoryCreateTestRole(container, registry)
			if err != nil {
				return nil, err
			}

			err = test.RepositoryCreateTestUserRole(container, registry)
			if err != nil {
				return nil, err
			}
		}
	}

	return service, nil
}
