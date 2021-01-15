package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestTruncateString_NotTruncated(t *testing.T) {
	str := "Hello world"
	assert.Equal(t, TruncateString(str, 50), "Hello world")
}

func TestTruncateString_Truncated(t *testing.T) {
	str := "The quick brown fox jumps over the lazy dog"
	maxLen := 25
	truncatedStr := TruncateString(str, maxLen)
	assert.Equal(t, truncatedStr, "The quick brown fox ju...")
	assert.Len(t, truncatedStr, maxLen)
}
