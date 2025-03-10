package management_test

import (
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/shared/model/management"
	"testing"
)

func TestVerifyCreateRolePayloadFailsOnEmptyName(t *testing.T) {
	payload := &management.RoleCreatePayload{
		Name: "",
	}

	err := payload.Verify()
	assert.Error(t, err, "verify should return an error")
	assert.Equal(t, err.Error(), "name is required")
}

func TestVerifyCreateRolePayloadSucceeds(t *testing.T) {
	payload := &management.RoleCreatePayload{
		Name: "admin",
	}

	err := payload.Verify()
	assert.NoError(t, err, "verify should not return an error")
}
