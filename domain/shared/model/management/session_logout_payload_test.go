package management_test

import (
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/shared/model/management"
	"testing"
)

func TestVerifySessionLogoutPayloadFailsOnMissingSessionIDsList(t *testing.T) {
	payload := &management.SessionLogoutPayload{
		SessionIDs: nil,
	}

	err := payload.Verify()
	assert.Error(t, err, "verify should return an error")
	assert.Equal(t, err.Error(), "session_ids is required")
}

func TestVerifySessionLogoutPayloadFailsOnEmptySessionIDsList(t *testing.T) {
	payload := &management.SessionLogoutPayload{
		SessionIDs: []string{},
	}

	err := payload.Verify()
	assert.Error(t, err, "verify should return an error")
	assert.Equal(t, err.Error(), "session_ids must contain at least one item")
}

func TestVerifySessionLogoutPayloadFailsOnEmptySessionID(t *testing.T) {
	payload := &management.SessionLogoutPayload{
		SessionIDs: []string{""},
	}

	err := payload.Verify()
	assert.Error(t, err, "verify should return an error")
	assert.Equal(t, err.Error(), "session_ids[0]: session id cannot be empty")
}

func TestVerifySessionLogoutPayloadSucceeds(t *testing.T) {
	payload := &management.SessionLogoutPayload{
		SessionIDs: []string{"session_id"},
	}

	err := payload.Verify()
	assert.NoError(t, err, "verify should not return an error")
}
