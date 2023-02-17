package eip4361

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const msg = `localhost:6006 wants you to sign in with your Ethereum account:
0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2

I accept the ServiceOrg Terms of Service: https://service.invalid/tos

URI: https://service.invalid/login
Version: 1
Chain ID: 1
Nonce: 32891756
Issued At: 2021-09-30T16:25:240Z
Request ID: F369349D-9B66-4367-BAF2-AE9D83E0F9FA
Resources:
- ipfs://bafybeiemxf5abjwjbikoz4mc3a3dla6ual3jsgpdr4cjr3oz3evfyavhwq/
- https://example.com/my-web2-claim.json`

func TestParse(t *testing.T) {
	m, err := Parse(msg)
	require.NoError(t, err)

	assert.Equal(t, "localhost:6006", m.Domain)
	assert.Equal(t, "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", m.Address)
	assert.Equal(t, "I accept the ServiceOrg Terms of Service: https://service.invalid/tos", m.Statement)
	assert.Equal(t, "https://service.invalid/login", m.URI)
	assert.Equal(t, "1", m.Version)
	assert.Equal(t, 1, m.ChainID)
	assert.Equal(t, "32891756", m.Nonce)
	assert.Equal(t, "2021-09-30T16:25:240Z", m.IssuedAt)
	assert.Empty(t, m.ExpirationTime)
	assert.Empty(t, m.NotBefore)
	assert.Equal(t, "F369349D-9B66-4367-BAF2-AE9D83E0F9FA", m.RequestID)
	assert.Equal(t, []string{
		"ipfs://bafybeiemxf5abjwjbikoz4mc3a3dla6ual3jsgpdr4cjr3oz3evfyavhwq/",
		"https://example.com/my-web2-claim.json",
	}, m.Resources)

	assert.Equal(t, msg, m.String())
}
