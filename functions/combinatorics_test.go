package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubfactorial(t *testing.T) {
	assert.Equal(t, Subfactorial(6), 265.0)
	assert.Equal(t, Subfactorial(10), 1334961.0)
	assert.Equal(t, Subfactorial(8), 14833.0)
}
