package eip4361

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Run("domain", func(t *testing.T) {
		assert.NoError(t, validateDomain("example.com"))
		assert.NoError(t, validateDomain("tom@example.com"))
		assert.NoError(t, validateDomain("tom@example.com:80"))
	})
}
