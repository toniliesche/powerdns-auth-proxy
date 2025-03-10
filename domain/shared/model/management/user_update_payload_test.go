package management_test

import (
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/shared/model/management"
	"testing"
)

func TestVerifyUserUpdatePayloadFailsOnMissingPassword(t *testing.T) {
	payload := management.UserUpdatePayload{}

	err := payload.Verify()
	assert.Error(t, err, "verify should return an error")
	assert.Equal(t, "password is required", err.Error())
}

func TestVerifyUserUpdatePayloadSucceeds(t *testing.T) {
	payload := management.UserUpdatePayload{
		Password: "password",
	}

	err := payload.Verify()
	assert.NoError(t, err, "verify should not return an error")
}
