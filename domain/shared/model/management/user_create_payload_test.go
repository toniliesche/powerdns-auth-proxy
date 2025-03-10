package management_test

import (
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/shared/model/management"
	"testing"
)

func TestVerifyUserCreatePayloadFailsOnEmptyPassword(t *testing.T) {
	payload := &management.UserCreatePayload{
		Username: "username",
		Password: "",
	}

	err := payload.Verify()
	assert.Error(t, err, "verify should return an error")
	assert.Equal(t, err.Error(), "password is required")
}

func TestVerifyUserCreatePayloadFailsOnEmptyUsername(t *testing.T) {
	payload := &management.UserCreatePayload{
		Username: "",
		Password: "password",
	}

	err := payload.Verify()
	assert.Error(t, err, "verify should return an error")
	assert.Equal(t, err.Error(), "username is required")
}

func TestVerifyUserCreatePayloadSucceeds(t *testing.T) {
	payload := &management.UserCreatePayload{
		Username: "username",
		Password: "password",
	}

	err := payload.Verify()
	assert.NoError(t, err, "verify should not return an error")
}
