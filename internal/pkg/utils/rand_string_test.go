package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandString_Length(t *testing.T) {
	assert.Len(t, RandString(5), 5)
}

func TestRandString_Randomness(t *testing.T) {
	assert.NotEqual(t, RandString(5), RandString(5))
	assert.NotEqual(t, RandString(4), RandString(4))
	assert.NotEqual(t, RandString(3), RandString(3))
	assert.NotEqual(t, RandString(2), RandString(2))
	assert.NotEqual(t, RandString(1), RandString(1))
}