package management_test

import (
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/shared/model/management"
	"testing"
)

func TestVerifyUserRolePayloadFailsOnMissingRolesList(t *testing.T) {
	payload := management.UserRolePayload{}

	err := payload.Verify()

	assert.Error(t, err, "verify should return an error")
	assert.Equal(t, "roles is required", err.Error())
}

func TestVerifyUserRolePayloadFailsOnEmptyRolesList(t *testing.T) {
	payload := management.UserRolePayload{
		Roles: []string{},
	}

	err := payload.Verify()

	assert.Error(t, err, "verify should return an error")
	assert.Equal(t, "roles must contain at least one item", err.Error())
}

func TestVerifyUserRolePayloadFailsOnEmptyRole(t *testing.T) {
	payload := management.UserRolePayload{
		Roles: []string{""},
	}

	err := payload.Verify()

	assert.Error(t, err, "verify should return an error")
	assert.Equal(t, "roles[0]: role cannot be empty", err.Error())
}

func TestVerifyUserRolePayloadSucceeds(t *testing.T) {
	payload := management.UserRolePayload{
		Roles: []string{"admin"},
	}

	err := payload.Verify()

	assert.NoError(t, err, "verify should not return an error")
}

func TestUserRolePayloadDetectsAdminAsAdminRole(t *testing.T) {
	payload := management.UserRolePayload{
		Roles: []string{"admin"},
	}

	assert.True(t, payload.ContainsAdminRole(), "admin role should be detected as admin")
}

func TestUserRolePayloadDetectsSuperAdminAsAdminRole(t *testing.T) {
	payload := management.UserRolePayload{
		Roles: []string{"superadmin"},
	}

	assert.True(t, payload.ContainsAdminRole(), "superadmin role should be detected as admin")
}

func TestUserRolePayloadDetectsAdminAndSuperAdminAsAdminRole(t *testing.T) {
	payload := management.UserRolePayload{
		Roles: []string{"admin", "superadmin"},
	}

	assert.True(t, payload.ContainsAdminRole(), "admin and superadmin roles should be detected as admin")
}

func TestUserRolePayloadDetectsNoAdminRole(t *testing.T) {
	payload := management.UserRolePayload{
		Roles: []string{"user"},
	}

	assert.False(t, payload.ContainsAdminRole(), "user role should not be detected as admin")
}
