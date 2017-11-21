package mhttp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrependHeaders(t *testing.T) {
	headers, err := PrependHeaders([]string{"X-CUSTOM: 1", "X-BLABLA: 2"})

	assert.Nil(t, err)
	assert.Equal(t, "1", headers["X-CUSTOM"])
	assert.Equal(t, "2", headers["X-BLABLA"])
}
