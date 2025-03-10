package management_test

import (
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/shared/model/management"
	"testing"
)

func TestVerifyDomainCreatePayloadFailsOnEmptyFqdn(t *testing.T) {
	payload := &management.DomainCreatePayload{
		Fqdn: "",
	}

	err := payload.Verify()
	assert.Error(t, err, "verify should return an error")
	assert.Equal(t, err.Error(), "fqdn is required")
}

func TestVerifyDomainCreatePayloadSucceeds(t *testing.T) {
	payload := &management.DomainCreatePayload{
		Fqdn: "example.com",
	}

	err := payload.Verify()
	assert.NoError(t, err, "verify should not return an error")
}
