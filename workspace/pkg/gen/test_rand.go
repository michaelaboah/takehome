package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomness(t *testing.T) {
	b1 := RandBytes(32)
	b2 := RandBytes(32)
	assert.NotEqual(t, b1, b2)

	s1 := RandString(32)
	s2 := RandString(32)
	assert.NotEqual(t, s1, s2)
}
