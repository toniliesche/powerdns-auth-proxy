package management_test

import (
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/shared/model/management"
	"testing"
)

func TestVerifyDomainRoleFailsOnMissingDomain(t *testing.T) {
	role := &management.DomainRole{
		Domain: "",
		Role:   "admin",
	}

	err := role.Verify()
	assert.Error(t, err, "verify should return an error")
	assert.Equal(t, err.Error(), "domain is required")
}

func TestVerifyDomainRoleFailsOnMissingRole(t *testing.T) {
	role := &management.DomainRole{
		Domain: "example.com",
		Role:   "",
	}

	err := role.Verify()
	assert.Error(t, err, "verify should return an error")
	assert.Equal(t, err.Error(), "role is required")
}

func TestVerifyDomainRoleSucceeds(t *testing.T) {
	role := &management.DomainRole{
		Domain: "example.com",
		Role:   "admin",
	}

	err := role.Verify()
	assert.NoError(t, err, "verify should not return an error")
}
